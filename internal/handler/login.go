package handler

import (
	"auth/internal/entity/user"
	"net/http"
	"github.com/labstack/echo/v4"
)

// LoginHandler godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.User true "User credentials"
// @Success 201 {object} map[string]interface{} "JWT token"
// @Failure 400 {string} string "invalid input"
// @Failure 500 {string} string "internal server error"
// @Router /login [post]
func (h *Handler) LoginHandler(c echo.Context) error {
	user := user.User{}

	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidInput)
	}

	token, err := h.service.Login(user.Username, user.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"token": token,
	})
}