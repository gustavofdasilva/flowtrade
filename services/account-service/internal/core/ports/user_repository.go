package ports

import (
	"account-service/internal/core/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByEmail(email string) (user *domain.User, err error)
	GetByID(id uuid.UUID) (user *domain.User, err error)
	CreateUser(user domain.User) (*domain.User, error)
	DeleteUserByID(id uuid.UUID) (err error)
	UpdateUser(user domain.User) (err error)
	IsEmailAlreadyInUse(email string, excludedID *uuid.UUID) (exists bool, err error)
	IsUsernameAlreadyInUse(username string, excludedID *uuid.UUID) (exists bool, err error)
}
