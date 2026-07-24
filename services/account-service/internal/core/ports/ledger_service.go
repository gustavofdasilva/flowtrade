package ports

import (
	"account-service/internal/core/domain"

	"github.com/google/uuid"
)

type LedgerService interface {
	GetStatement(userID uuid.UUID, accountID uuid.UUID, page, limit int) ([]*domain.LedgerEntry, int, error)
}
