package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tsntt/footballapi/internal/dto"
	"github.com/tsntt/footballapi/pkg/utils"
)

type AuthMiddleware struct {
	jwtService *utils.JWTService
}

func NewAuthMiddleware(jwtService *utils.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// JWTAuth middleware
func (m *AuthMiddleware) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				slog.Error("Authorization header required")
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
			}

			// Check token format "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				slog.Error("Invalid authorization header format")
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
			}

			tokenString := parts[1]
			claims, err := m.jwtService.ValidateToken(tokenString)
			if err != nil {
				slog.Error("Invalid token", slog.String("err", err.Error()))
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
			}

			// Set claims into context
			c.Set("user", claims)
			return next(c)
		}
	}
}

// Check if is admin
func (m *AuthMiddleware) AdminAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*dto.JWTClaims)
			if !ok {
				slog.Error("User not found in context")
				return echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
			}

			if user.Role != "admin" {
				slog.Error("Admin access required")
				return echo.NewHTTPError(http.StatusForbidden, "Admin access required")
			}

			return next(c)
		}
	}
}

// Get user from context
func GetUserFromContext(c echo.Context) (*dto.JWTClaims, error) {
	user, ok := c.Get("user").(*dto.JWTClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
	}
	return user, nil
}
