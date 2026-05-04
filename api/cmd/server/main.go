package main

import (
	"log/slog"

	"github.com/RintaroNasu/drive-route-planner/api/internal/httpx"
	"github.com/RintaroNasu/drive-route-planner/api/internal/logging"
	"github.com/RintaroNasu/drive-route-planner/api/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	logger := logging.New()
	slog.SetDefault(logger)

	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)
	routes.NewRouter(e)

	e.Logger.Fatal(e.Start(":8080"))
}
