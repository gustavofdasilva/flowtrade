package ports

import (
	"account-service/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

type LedgerService interface {
	GetStatement(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, page, limit int) ([]domain.LedgerEntry, int, error)
}
