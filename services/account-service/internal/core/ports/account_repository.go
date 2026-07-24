package ports

import (
	"account-service/internal/core/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	GetByID(id uuid.UUID, userID uuid.UUID) (*domain.Account, error)
	Create(account *domain.Account) (*domain.Account, error)
	AddBalance(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	SubtractBalance(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error)
	// CreateLedgerEntry(entry *domain.LedgerEntry) error
	// GetStatement(accountID uuid.UUID, page, limit int) ([]*domain.LedgerEntry, error)
}
