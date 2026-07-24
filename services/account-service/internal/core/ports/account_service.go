package ports

import (
	"account-service/internal/core/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountService interface {
	Get(id uuid.UUID, userID uuid.UUID) (*domain.Account, error)
	Create(userID uuid.UUID, currency string) (*domain.Account, error)
	Deposit(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	Withdrawal(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	CheckBalance(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (bool, error)
}
