package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tsntt/footballapi/internal/config"
)

// CustomValidator para o Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Retorna 400 Bad Request
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	// Carregar Variáveis de Ambiente
	if err := os.Setenv("PORT", config.LoadConfig().Port); err != nil {
		log.Fatal("Erro ao setar variável PORT:", err)
	}

	cfg := config.LoadConfig()

	// 1. Conexão com o Banco de Dados (Postgres)
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// 2. Inicializar Repositórios e Broadcaster
	userRepo := repository.NewUserRepository(db)
	fanRepo := repository.NewFanRepository(db)

	// Broadcaster (Sistema de mensagens com Channels)
	bcast := broadcaster.NewBroadcaster()
	go bcast.Start() // Inicia a goroutine do Broadcaster

	// 3. Inicializar Handlers
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret)
	championshipHandler := handler.NewChampionshipHandler(cfg.FootballAPIKey)
	fanHandler := handler.NewFanHandler(fanRepo)
	adminHandler := handler.NewAdminHandler(db, bcast) // db para listar partidas, bcast para o broadcast

	// 4. Configurar Echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 5. Configurar Middlewares e Rotas

	// Middleware de Autorização JWT
	authMiddleware := auth.NewAuthMiddleware(cfg.JWTSecret)

	// Rotas de Autenticação
	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Rotas Públicas de Campeonato
	e.GET("/championship", championshipHandler.ListChampionships)
	e.GET("/championship/:id/matches", championshipHandler.ListMatches)

	// Rotas Protegidas (requerem JWT)
	protectedGroup := e.Group("")
	protectedGroup.Use(authMiddleware.ValidateJWT)

	// Rota de Fã
	protectedGroup.POST("/fans", fanHandler.SubscribeFan)

	// Rotas de Administrador (você pode adicionar uma verificação de role aqui)
	adminGroup := protectedGroup.Group("/admin")
	adminGroup.GET("", adminHandler.ListAdminMatches)                            // Listar partidas com estados
	adminGroup.POST("/broadcast/:match_id", adminHandler.BroadcastMatchUpdate)   // Disparar Broadcast
	adminGroup.POST("/broadcast/cancel/:match_id", adminHandler.CancelBroadcast) // Cancelar Broadcast

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Printf("Servidor rodando na porta %s", cfg.Port)
	go func() {
		if err := e.StartServer(s); err != nil && err != http.ErrServerClosed {
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
