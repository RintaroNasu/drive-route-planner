package routes

import (
	"github.com/RintaroNasu/drive-route-planner/api/internal/handler"
	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo) {
	h := handler.NewRouteHandler()

	e.POST("/route-from-place", h.RouteFromPlace)
}
