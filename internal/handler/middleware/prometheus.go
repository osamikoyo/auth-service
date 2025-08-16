package middleware

import (
	"auth/internal/metrics"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			duration := time.Since(start).Seconds()

			status := c.Response().Status

			metrics.RequestCount.WithLabelValues(
				c.Request().Method,
				c.Request().URL.Path,
				http.StatusText(status),
			)

			metrics.RequestDuration.WithLabelValues(
				c.Request().Method,
				c.Request().URL.Path,
				http.StatusText(status),
			).Observe(duration)

			return err
		}
	}
}