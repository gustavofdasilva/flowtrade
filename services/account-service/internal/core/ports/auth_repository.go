package ports

import (
	"account-service/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type AuthRepository interface {
	SaveRefreshToken(userID uuid.UUID, refreshToken string, expiresAt time.Time) error
	GetRefreshToken(refreshToken string) (*domain.RefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
	RevokeAllRefreshTokens(userID uuid.UUID) error
}
