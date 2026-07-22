package ports

import (
	"account-service/internal/core/domain"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(email string, password string) (*domain.AuthResult, error)
	Refresh(refreshToken string) (*domain.AuthResult, error)
	Logout(refreshToken string) error
	LogoutAll(userID uuid.UUID) error
}
