package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	data "github.com/tsntt/footballapi/data/postgres"
	"github.com/tsntt/footballapi/internal/api/handler"
	"github.com/tsntt/footballapi/internal/api/middleware"
	"github.com/tsntt/footballapi/internal/config"
	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/pkg/broadcast"
	consumer "github.com/tsntt/footballapi/pkg/external_api_consumer"
	"github.com/tsntt/footballapi/pkg/services/email"
	"github.com/tsntt/footballapi/pkg/services/sms"
	"github.com/tsntt/footballapi/pkg/utils"

	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// load config
	cfg := config.Load()

	// Connect to database
	db, err := data.NewDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
		cfg.Database.Port,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// init repositories
	userRepo := data.NewUserRepository(db)
	fanRepo := data.NewFanRepository(db)
	broadcastRepo := data.NewBroadcastRepository(db)

	// init services
	jwtService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresHours)
	footballAPI := consumer.NewFootballAPIClient(cfg.FootballAPI.URL, cfg.FootballAPI.Token)
	emailService := email.NewMailgunService(cfg.Server.Host, cfg.EmailAPI.APIKey, cfg.EmailAPI.From)
	smsService := sms.NewTwilioService(cfg.SMSAPI.AccountSID, cfg.SMSAPI.APIKey, cfg.SMSAPI.From, "")
	broadcastService := broadcast.NewBroadcastService()

	broadcastService.RegisterNotifier(broadcast.Email, emailService)
	broadcastService.RegisterNotifier(broadcast.SMS, smsService)

	// init controllers
	userController := controller.NewUserController(userRepo, jwtService)
	championshipController := controller.NewChampionshipController(footballAPI)
	fanController := controller.NewFanController(fanRepo)
	adminController := controller.NewAdminController(
		footballAPI,
		fanRepo,
		broadcastRepo,
		broadcastService,
	)

	// init handlers
	handlers := handler.NewHandlers(
		userController,
		championshipController,
		fanController,
		adminController,
	)

	// init middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Configure Echo
	e := echo.New()

	// Middlewares globais
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Configure rotas
	handler.SetupRoutes(e, handlers, authMiddleware)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// graceful shutdown
	port := ":" + cfg.Server.Port

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	slog.Info("Starting server on port", slog.String("port", cfg.Server.Port))

	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
