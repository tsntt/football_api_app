package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/tsntt/footballapi/internal/controller"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type AdminHandler struct {
	controller *controller.AdminController
}

func NewAdminHandler(controller *controller.AdminController) *AdminHandler {
	return &AdminHandler{controller: controller}
}

func (h *AdminHandler) GetMatches(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer ws.Close()

	defer h.controller.UnregisterWS(ws)

	fmt.Println("New admin connected")

	for {
		if _, _, err := ws.NextReader(); err != nil {
			fmt.Println("Admin connection lost!")
			break
		}
	}

	matches, err := h.controller.GetMatches(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: maybe use a websocket client instead
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
