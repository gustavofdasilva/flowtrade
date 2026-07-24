package ports

import (
	"account-service/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

type LedgerRepository interface {
	CreateLedgerEntry(ctx context.Context, entry *domain.LedgerEntry) error
	GetStatement(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, limit, offset int) ([]domain.LedgerEntry, int, error)
}
