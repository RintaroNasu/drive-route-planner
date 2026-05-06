package httpx

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(l *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		ctx := c.Request().Context()

		if ae, ok := err.(*AppError); ok {

			// ログ分岐
			if ae.Status >= 500 {
				l.ErrorContext(ctx, "server_error",
					"code", ae.Code,
					"error", ae.Err,
					"path", c.Path(),
					"method", c.Request().Method,
				)
			} else {
				l.WarnContext(ctx, "client_error",
					"code", ae.Code,
					"error", ae.Err,
					"path", c.Path(),
					"method", c.Request().Method,
				)
			}
			_ = c.JSON(ae.Status, map[string]interface{}{
				"error": map[string]string{
					"code":    ae.Code,
					"message": ae.Message,
				},
			})
			return
		}

		// 想定外
		l.ErrorContext(ctx, "unexpected_error",
			"error", err,
			"path", c.Path(),
			"method", c.Request().Method,
		)

		_ = c.JSON(500, map[string]interface{}{
			"error": map[string]string{
				"code":    "InternalError",
				"message": "internal server error",
			},
		})
		return
	}
}
