package ports

import (
	"account-service/internal/core/domain"

	"github.com/google/uuid"
)

type UserService interface {
	Register(user domain.User) (*domain.User, error)
	Update(user domain.User) (*domain.User, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.User, error)
}
