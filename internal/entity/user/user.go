package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UID       string `json:"uuid" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"primaryKey"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(username, password string) *User {
	return &User{
		UID:       uuid.New().String(),
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
	}
}
