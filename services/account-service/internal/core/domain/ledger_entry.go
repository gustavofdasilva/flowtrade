package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LedgerType string

const (
	LedgerTypeDeposit     = "DEPOSIT"
	LedgerTypeWithdrawal  = "WITHDRAWAL"
	LedgerTypeTradeDebit  = "TRADE_DEBIT"
	LedgerTypeTradeCredit = "TRADE_CREDIT"
)

type LedgerEntry struct {
	ID           uuid.UUID
	AccountID    uuid.UUID
	Type         LedgerType
	Amount       decimal.Decimal
	BalanceAfter decimal.Decimal
	Description  string
	ReferenceID  *uuid.UUID
	CreatedAt    time.Time
}
