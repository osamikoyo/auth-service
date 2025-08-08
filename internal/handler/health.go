package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Accept json
// @Produce plain
// @Success 200 {string} string "ok"
// @Router /health [get]
func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}