package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUserHandler godoc
// @Summary Get user account
// @Description Retrieve user account details by UID (requires authentication)
// @Tags account
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} user.User "User account details"
// @Failure 400 {string} string "failed get uid"
// @Failure 500 {string} string "internal server error"
// @Router /account/get [get]
func (h *Handler) GetUserHandler(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok {
		return c.String(http.StatusBadRequest, "failed get uid")
	}

	account, err := h.service.GetAccount(uid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}