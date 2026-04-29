package handler

import (
	"net/http"

	"github.com/RintaroNasu/drive-route-planner/api/internal/models"
	"github.com/RintaroNasu/drive-route-planner/api/internal/service"
	"github.com/labstack/echo/v4"
)

type RouteHandler struct {
	service *service.RouteService
}

func NewRouteHandler() *RouteHandler {
	return &RouteHandler{
		service: service.NewRouteService(),
	}
}

func (h *RouteHandler) RouteFromPlace(c echo.Context) error {
	var req models.RouteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	result, err := h.service.RouteFromPlace(req.Places)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}
