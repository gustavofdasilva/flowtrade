package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Currency string

const (
	CurrencyBRL = "BRL"
	CurrencyUSD = "USD"
)

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Balance   decimal.Decimal
	Currency  Currency
	CreatedAt time.Time
	UpdatedAt time.Time
}
