package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoggingMV() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_custom}, host=${host},  method=${method}," +
			"uri=${uri}, status=${status}, latency=${latency_human}\n",
		CustomTimeFormat: "02.01.2006 15:04:05",
		Output:           os.Stdout,
	})
}

func Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}
