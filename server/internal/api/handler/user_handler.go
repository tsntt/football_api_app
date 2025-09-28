package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/dto"
)

type UserHandler struct {
	controller *controller.UserController
}

func NewUserHandler(controller *controller.UserController) *UserHandler {
	return &UserHandler{controller: controller}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req dto.UserRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.controller.Register(c.Request().Context(), &req)
	if err != nil {
		slog.Error("Failed to register user", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c echo.Context) error {
	var req dto.UserRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.controller.Login(c.Request().Context(), &req)
	if err != nil {
		slog.Error("Failed to login user", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// INFO: would be safer to store token in cookie
	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Logout(c echo.Context) error {
	// INFO: would be safer to remove token from cookie
	response, err := h.controller.Logout(c.Request().Context())
	if err != nil {
		slog.Error("Failed to logout user", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
