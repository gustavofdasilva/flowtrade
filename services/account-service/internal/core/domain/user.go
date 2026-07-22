package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleConsumer UserRole = "CONSUMER"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      UserRole  `json:"role"`
	CreateAt  time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
