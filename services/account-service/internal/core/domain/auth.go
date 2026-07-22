package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         *User
}
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (r *RefreshToken) IsValid() bool {
	if !r.RevokedAt.IsZero() || time.Now().After(r.ExpiresAt) {
		return false
	}

	return true
}
