package ports

import (
	"account-service/internal/core/domain"
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountService interface {
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Account, error)
	Create(ctx context.Context, userID uuid.UUID, currency domain.Currency) (*domain.Account, error)
	Deposit(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	Withdrawal(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	CheckBalance(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (bool, error)
}
