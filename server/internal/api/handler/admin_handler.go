package handler

import (
	"log/slog"
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
	matches, err := h.controller.GetMatches(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, matches)
}

func (h *AdminHandler) WsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		slog.Error("Failed to upgrade to websocket", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer ws.Close()

	h.controller.RegisterWS(ws)
	defer h.controller.UnregisterWS(ws)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			slog.Info("read:", slog.String("err", err.Error()))
			break
		}
		slog.Info("recv:", slog.String("msg", string(msg)))

		if _, _, err := ws.NextReader(); err != nil {
			slog.Info("Admin connection lost!")
			break
		}
	}

	return nil
}

func (h *AdminHandler) BroadcastMatch(c echo.Context) error {
	matchIDstr := c.Param("match_id")

	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		slog.Error("Invalid match ID", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid match ID")
	}

	response, err := h.controller.BroadcastMatch(c.Request().Context(), matchID)
	if err != nil {
		slog.Error("Failed to broadcast match", slog.String("err", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
