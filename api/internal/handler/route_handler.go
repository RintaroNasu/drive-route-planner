package handler

import (
	"net/http"

	"github.com/RintaroNasu/drive-route-planner/api/internal/httpx"
	"github.com/RintaroNasu/drive-route-planner/api/internal/models"
	"github.com/RintaroNasu/drive-route-planner/api/internal/service"
	"github.com/labstack/echo/v4"
)

type RouteHandler struct {
	service service.RouteService
}

func NewRouteHandler(s service.RouteService) *RouteHandler {
	return &RouteHandler{
		service: s,
	}
}

func (h *RouteHandler) RouteFromPlace(c echo.Context) error {
	var req models.RouteRequest

	if err := c.Bind(&req); err != nil {
		return httpx.InvalidRequest("invalid request", err)
	}

	result, err := h.service.RouteFromPlace(req.Places)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
