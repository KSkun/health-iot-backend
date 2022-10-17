package middleware

import (
	"github.com/KSkun/health-iot-backend/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"strings"
)

func Logger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[Echo] ${time_custom} ${status} ${method} ${uri} ${latency_human} ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Skipper: func(ctx echo.Context) bool {
			log.Printf("[Echo] User %s logged in with status %d", ctx.QueryParam("name"), ctx.Response().Status)
			return !config.Debug && strings.Contains(ctx.Request().RequestURI, "token")
		},
	})
}
