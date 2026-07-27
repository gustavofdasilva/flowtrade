package services

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AssetService struct {
	repo ports.AssetRepository
}

func NewAssetService(repo ports.AssetRepository) *AssetService {
	return &AssetService{
		repo: repo,
	}
}

func (s *AssetService) GetAll(ctx context.Context, page, limit int) ([]domain.Asset, int, error) {
	offset := (page - 1) * limit

	return s.repo.GetAll(ctx, offset, limit)

}

func (s *AssetService) GetByTicker(ctx context.Context, ticker string) (*domain.Asset, error) {
	return s.repo.GetByTicker(ctx, ticker)
}

func (s *AssetService) GetPricesByTicker(ctx context.Context, ticker string) ([]domain.AssetPrice, error) {
	return s.repo.GetPricesByTicker(ctx, ticker)
}

func (s *AssetService) CreateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	return s.repo.CreateAsset(ctx, asset)
}

func (s *AssetService) UpdateAssetPriceByTicker(ctx context.Context, source, ticker string, price decimal.Decimal) (*domain.Asset, error) {
	asset, err := s.repo.GetByTicker(ctx, ticker)
	if err != nil {
		return nil, err
	}

	assetPrice := domain.AssetPrice{
		AssetID: asset.ID,
		Price:   price,
		Source:  source,
	}

	newAssetPrice, err := s.repo.CreateAssetPrice(ctx, assetPrice)

	asset.Price = newAssetPrice.Price

	return asset, nil
}

func (s *AssetService) UpdateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	return s.repo.UpdateAsset(ctx, asset)
}

func (s *AssetService) Delete(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	return s.repo.Delete(ctx, id)
}
