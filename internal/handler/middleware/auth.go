package middleware

import (
	"auth/pkg/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	ErrAuth = "authentification error"
)

type Auth struct {
	logger *logger.Logger
	key    string
}

func NewAuth(key string, logger *logger.Logger) *Auth {
	return &Auth{
		key:    key,
		logger: logger,
	}
}

func (a *Auth) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.key, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *Auth) GetUserIDFromToken(tokenString string) (string, error) {
	token, err := a.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userID, ok := claims["uid"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}

func (a *Auth) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		a.logger.Info("new auth request", zap.String("token", token))

		if len(token) == 0 {
			a.logger.Error("token is empty", zap.String("path", c.Path()))

			return c.String(http.StatusBadGateway, ErrAuth)
		}

		tokenStr := strings.TrimPrefix(token, "Bearer ")
		if token == tokenStr {
			a.logger.Error("bearer required")

			return c.String(http.StatusBadGateway, "bearer requiered")
		}

		uid, err := a.GetUserIDFromToken(tokenStr)
		if err != nil {
			a.logger.Error("failed get user id from token", zap.String("token", token))

			return c.String(http.StatusBadGateway, ErrAuth)
		}

		c.Set("uid", uid)

		return next(c)
	}
}
