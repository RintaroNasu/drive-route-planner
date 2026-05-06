package routes

import (
	"github.com/RintaroNasu/drive-route-planner/api/internal/handler"
	"github.com/RintaroNasu/drive-route-planner/api/internal/service"
	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo) {
	svc := service.NewRouteService()
	h := handler.NewRouteHandler(svc)

	e.POST("/route-from-place", h.RouteFromPlace)
}
