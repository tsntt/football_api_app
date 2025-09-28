package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/tsntt/footballapi/internal/api/middleware"
	"github.com/tsntt/footballapi/internal/controller"
)

type Handlers struct {
	User         *UserHandler
	Championship *ChampionshipHandler
	Fan          *FanHandler
	Admin        *AdminHandler
}

func NewHandlers(
	userController *controller.UserController,
	championshipController *controller.ChampionshipController,
	fanController *controller.FanController,
	adminController *controller.AdminController,
) *Handlers {
	return &Handlers{
		User:         NewUserHandler(userController),
		Championship: NewChampionshipHandler(championshipController),
		Fan:          NewFanHandler(fanController),
		Admin:        NewAdminHandler(adminController),
	}
}

func SetupRoutes(e *echo.Echo, handlers *Handlers, authMiddleware *middleware.AuthMiddleware) {
	apiV1 := e.Group("/api/v1")
	// Public
	auth := apiV1.Group("/auth")
	auth.POST("/register", handlers.User.Register)
	auth.POST("/login", handlers.User.Login)
	auth.POST("/logout", handlers.User.Logout)

	// Protected
	protected := apiV1.Group("")
	protected.Use(authMiddleware.JWTAuth())

	// Championship
	protected.GET("/championship", handlers.Championship.GetChampionships)
	protected.GET("/championship/:id/matches", handlers.Championship.GetMatches)

	// Fan
	protected.POST("/fans", handlers.Fan.Subscribe)
	protected.DELETE("/fans", handlers.Fan.Unsubscribe)
	protected.GET("/fans", handlers.Fan.GetSubscriptions)

	// Protected [Only Admin]
	admin := apiV1.Group("/admin")
	admin.Use(authMiddleware.JWTAuth())
	admin.Use(authMiddleware.AdminAuth())
	admin.GET("/", handlers.Admin.GetMatches)
	admin.POST("/broadcast/:match_id", handlers.Admin.BroadcastMatch)
}
