package routes

import (
	"github.com/RintaroNasu/drive-route-planner/api/internal/handler"
	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	h := handler.NewRouteHandler()

	e.POST("/route-from-place", h.RouteFromPlace)

	return e
}
