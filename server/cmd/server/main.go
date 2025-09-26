package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	data "github.com/tsntt/footballapi/data/postgres"
	"github.com/tsntt/footballapi/internal/api/handler"
	"github.com/tsntt/footballapi/internal/api/middleware"
	"github.com/tsntt/footballapi/internal/config"
	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/pkg/broadcaster"
	consumer "github.com/tsntt/footballapi/pkg/external_api_consumer"
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
	broadcastService := broadcaster.NewBroadcastService(5)

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
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Shutting down server...")

		broadcastService.Stop()

		if err := e.Shutdown(nil); err != nil {
			log.Fatal("Server shutdown error:", err)
		}
	}()

	// start server
	port := ":" + cfg.Server.Port
	log.Printf("Starting server on port %s", cfg.Server.Port)

	if err := e.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
