package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// DeleteUser godoc
// @Summary Delete user account
// @Description Delete a user account by UID (requires authentication)
// @Tags account
// @Accept json
// @Produce plain
// @Security ApiKeyAuth
// @Success 200 {string} string "deleted"
// @Failure 400 {string} string "failed get uid"
// @Failure 500 {string} string "internal server error"
// @Router /account/delete [delete]
func (h *Handler) DeleteUser(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return c.String(http.StatusBadRequest, "failed get uid")
	}

	if err := h.service.DeleteAccount(uid); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "deleted")
}