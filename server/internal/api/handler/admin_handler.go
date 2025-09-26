package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tsntt/footballapi/internal/controller"
)

type AdminHandler struct {
	controller *controller.AdminController
}

func NewAdminHandler(controller *controller.AdminController) *AdminHandler {
	return &AdminHandler{controller: controller}
}

func (h *AdminHandler) GetMatches(c echo.Context) error {
	matches, err := h.controller.GetMatches(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, matches)
}

func (h *AdminHandler) BroadcastMatch(c echo.Context) error {
	matchIDstr := c.Param("match_id")

	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid match ID")
	}

	response, err := h.controller.BroadcastMatch(c.Request().Context(), matchID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
