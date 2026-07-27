package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AssetType string

const (
	AssetTypeStock  AssetType = "STOCK"
	AssetTypeETF    AssetType = "ETF"
	AssetTypeCrypto AssetType = "CRYPTO"
)

type AssetPriceFilter struct {
	From     *time.Time
	To       *time.Time
	Interval *time.Duration
}

type Asset struct {
	ID             uuid.UUID
	Ticker         string
	Name           string
	Type           AssetType
	Currency       Currency
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Price          decimal.Decimal
	PriceUpdatedAt time.Time
}

type AssetPrice struct {
	ID        uuid.UUID
	AssetID   uuid.UUID
	Price     decimal.Decimal
	Source    string
	CreatedAt time.Time
}
