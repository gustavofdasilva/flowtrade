package ports

import (
	"account-service/internal/core/domain"
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AssetService interface {
	GetAll(ctx context.Context, page, limit int) ([]domain.Asset, int, error)
	GetByTicker(ctx context.Context, ticker string) (*domain.Asset, error)
	GetPricesByTicker(ctx context.Context, ticker string) ([]domain.AssetPrice, error)
	CreateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error)
	UpdateAssetPriceByTicker(ctx context.Context, source, ticker string, price decimal.Decimal) (*domain.Asset, error)
	UpdateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error)
	Delete(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
}
