package handler

import (
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	//Use user_id from JWT security reasons
	req.UserID = user.UserID

	response, err := h.controller.Subscribe(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
