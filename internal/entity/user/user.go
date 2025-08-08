package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	UID       string    `json:"uuid" gorm:"primaryKey" swag:"description=Unique identifier of the user"`
	Username  string    `json:"username" gorm:"unique" swag:"description=Username for authentication"`
	Password  string    `json:"password" swag:"description=Password for authentication"`
	CreatedAt time.Time `json:"created_at" swag:"description=Timestamp of user creation"`
}

func NewUser(username, password string) *User {
	return &User{
		UID:       uuid.New().String(),
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
	}
}
