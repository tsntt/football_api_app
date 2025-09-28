package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tsntt/footballapi/internal/api/middleware"
	"github.com/tsntt/footballapi/internal/controller"
	"github.com/tsntt/footballapi/internal/dto"
)

type FanHandler struct {
	controller *controller.FanController
}

func NewFanHandler(controller *controller.FanController) *FanHandler {
	return &FanHandler{controller: controller}
}

func (h *FanHandler) Subscribe(c echo.Context) error {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req dto.FanRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	req.UserID = user.UserID

	response, err := h.controller.Subscribe(c.Request().Context(), &req)
	if err != nil {
		slog.Error("Failed to subscribe to team", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *FanHandler) Unsubscribe(c echo.Context) error {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req dto.UnsubscribeRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Invalid request body", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.controller.Unsubscribe(c.Request().Context(), user.UserID, &req)
	if err != nil {
		slog.Error("Failed to unsubscribe from team", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *FanHandler) GetSubscriptions(c echo.Context) error {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	subscriptions, err := h.controller.GetSubscriptions(c.Request().Context(), user.UserID)
	if err != nil {
		slog.Error("Failed to get user subscriptions", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, subscriptions)
}
