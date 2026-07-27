package ports

import (
	"account-service/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

type AssetRepository interface {
	GetAll(ctx context.Context, offset, limit int, typeFilter []domain.AssetType) ([]domain.Asset, int, error)
	GetByTicker(ctx context.Context, ticker string) (*domain.Asset, error)
	GetPricesByTicker(ctx context.Context, ticker string, filter domain.AssetPriceFilter) ([]domain.AssetPrice, error)
	CreateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error)
	CreateAssetPrice(ctx context.Context, asset domain.AssetPrice) (*domain.AssetPrice, error)
	UpdateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error)
	Delete(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
}
