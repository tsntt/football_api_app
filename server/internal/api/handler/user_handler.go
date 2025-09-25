package handler

import (
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.controller.Register(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c echo.Context) error {
	var req dto.UserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.controller.Login(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
