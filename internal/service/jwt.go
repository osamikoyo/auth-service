package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJwtKey(uid, key string, dur time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid" : uid,
		"exp" : time.Now().Add(dur).Unix(),
		"iat" : time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}