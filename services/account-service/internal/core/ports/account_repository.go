package ports

import (
	"account-service/internal/core/domain"
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Account, error)
	Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
	AddBalance(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	SubtractBalance(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
}
