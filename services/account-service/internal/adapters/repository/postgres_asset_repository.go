package repository

import (
	"account-service/internal/core/domain"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type PostgresAssetRepository struct {
	db *sql.DB
}

func NewPostgresAssetRepository(db *sql.DB) *PostgresAssetRepository {
	return &PostgresAssetRepository{db: db}
}

func (repo *PostgresAssetRepository) GetAll(ctx context.Context, offset, limit int) ([]domain.Asset, int, error) {
	query := `
		SELECT 
			id, 
			ticker,
			name,
			type,
			currency,
			is_active,
			created_at,
			updated_at,
			price, 
			price_updated_at,
			COUNT(*) OVER() as total
			FROM (
			SELECT DISTINCT ON (a.id)
				a.id, 
				a.ticker,
				a.name,
				a.type,
				a.currency,
				a.is_active,
				a.created_at,
				a.updated_at,
				ap.price, 
				ap.created_at price_updated_at,
				COUNT(a.id) OVER() as total 
			FROM assets a
			LEFT JOIN asset_prices ap ON a.id = ap.asset_id
			ORDER BY a.id, ap.created_at DESC
			LIMIT $1 OFFSET $2
		);	
	`

	rows, err := repo.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var assets []domain.Asset
	var total int
	for rows.Next() {
		var asset domain.Asset
		err = rows.Scan(
			&asset.ID,
			&asset.Ticker,
			&asset.Name,
			&asset.Type,
			&asset.Currency,
			&asset.Active,
			&asset.CreatedAt,
			&asset.UpdatedAt,
			&asset.Price,
			&asset.PriceUpdatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, err
		}

		assets = append(assets, asset)
	}

	return assets, total, nil
}

func (repo *PostgresAssetRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Asset, error) {
	query := `
		SELECT DISTINCT ON (a.id)
			a.id, 
			a.ticker,
			a.name,
			a.type,
			a.currency,
			a.is_active,
			a.created_at,
			a.updated_at,
			ap.price, 
			ap.created_at price_updated_at,
			COUNT(*) OVER() as total 
		FROM assets a
		LEFT JOIN asset_prices ap ON a.id = ap.asset_id
		WHERE a.ticker = $1
		ORDER BY a.id, ap.created_at DESC;
	`

	row := repo.db.QueryRowContext(ctx, query, ticker)

	var asset domain.Asset
	err := row.Scan(
		&asset.ID,
		&asset.Ticker,
		&asset.Name,
		&asset.Type,
		&asset.Currency,
		&asset.Active,
		&asset.CreatedAt,
		&asset.UpdatedAt,
		&asset.Price,
		&asset.PriceUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (repo *PostgresAssetRepository) GetPricesByTicker(ctx context.Context, ticker string) ([]domain.AssetPrice, error) {
	query := `
		SELECT 
			id,
			price,
			"source",
			created_at 
		FROM asset_prices ap 
		WHERE ap.asset_id = (SELECT id FROM assets WHERE ticker = $1)
		ORDER BY created_at DESC;
	`

	rows, err := repo.db.QueryContext(ctx, query, ticker)
	if err != nil {
		return nil, err
	}

	var assetPrices []domain.AssetPrice
	for rows.Next() {
		var price domain.AssetPrice
		err = rows.Scan(
			&price.ID,
			&price.Price,
			&price.Source,
			&price.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		assetPrices = append(assetPrices, price)
	}

	return assetPrices, nil
}

func (repo *PostgresAssetRepository) CreateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	query := `
		INSERT INTO assets
			(ticker, name, type, currency, is_active, created_at, updated_at)
		VALUES
			($1,$2,$3,$4,TRUE,now(),now())
		RETURNING
			id, ticker, name, type, currency, is_active, created_at, updated_at
	`

	err := repo.db.QueryRowContext(ctx, query, asset.Ticker, asset.Type, asset.Currency).Scan(
		&asset.ID,
		&asset.Ticker,
		&asset.Name,
		&asset.Type,
		&asset.Currency,
		&asset.Active,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (repo *PostgresAssetRepository) CreateAssetPrice(ctx context.Context, assetPrice domain.AssetPrice) (*domain.AssetPrice, error) {
	query := `
		INSERT INTO asset_prices
			(asset_id, price, source, created_at)
		VALUES
			($1,$2,$3,now())
		RETURNING
			id, asset_id, price, source, created_at
	`

	err := repo.db.QueryRowContext(ctx, query, assetPrice.ID, assetPrice.AssetID, assetPrice.Price, assetPrice.Source, assetPrice.CreatedAt).Scan(
		&assetPrice.ID,
		&assetPrice.AssetID,
		&assetPrice.Price,
		&assetPrice.Source,
		&assetPrice.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &assetPrice, nil
}

func (repo *PostgresAssetRepository) UpdateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	query := `
		UPDATE assets SET
			ticker=$1,
			name=$2,
			type=$3,
			currency=$4,
			is_active=$5,
			updated_at=$6
		WHERE
			id=$6
		RETURNING
			id, ticker, name, type, currency, is_active, created_at, updated_at

	`

	err := repo.db.QueryRowContext(ctx, query, asset.Ticker, asset.Type, asset.Currency, asset.Active, asset.ID).Scan(
		&asset.ID,
		&asset.Ticker,
		&asset.Name,
		&asset.Type,
		&asset.Currency,
		&asset.Active,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (repo *PostgresAssetRepository) Delete(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	query := `
		DELETE FROM assets WHERE id=$1
		RETURNING
			id, ticker, name, type, currency, is_active, created_at, updated_at
	`

	var asset domain.Asset
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&asset.ID,
		&asset.Ticker,
		&asset.Name,
		&asset.Type,
		&asset.Currency,
		&asset.Active,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}
