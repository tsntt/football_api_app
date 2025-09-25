package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tsntt/footballapi/internal/controller"
)

type ChampionshipHandler struct {
	controller *controller.ChampionshipController
}

func NewChampionshipHandler(controller *controller.ChampionshipController) *ChampionshipHandler {
	return &ChampionshipHandler{controller: controller}
}

func (h *ChampionshipHandler) GetChampionships(c echo.Context) error {
	championships, err := h.controller.GetChampionships(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, championships)
}

func (h *ChampionshipHandler) GetMatches(c echo.Context) error {
	championshipID := c.Param("id")
	team := c.QueryParam("team")
	stage := c.QueryParam("stage")

	matches, err := h.controller.GetMatches(c.Request().Context(), championshipID, team, stage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, matches)
}
