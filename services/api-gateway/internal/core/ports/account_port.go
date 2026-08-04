package ports

import (
	"api-gateway/internal/core/domain"
	"context"
)

type AccountClient interface {
	Register(ctx context.Context, input domain.RegisterInput) error
	Login(ctx context.Context, input domain.LoginInput) (*domain.AuthOutput, error)
	RefreshToken(ctx context.Context, input domain.RefreshTokenInput) (*domain.AuthOutput, error)
	Logout(ctx context.Context, input domain.LogoutInput) error
	LogoutAll(ctx context.Context, input domain.UserIDInput) error
	UpdateUser(ctx context.Context, input domain.UpdateUserInput) (*domain.UserOutput, error)
	DeleteUser(ctx context.Context, input domain.UserIDInput) error
	GetMe(ctx context.Context, input domain.UserIDInput) (*domain.UserOutput, error)

	GetStatement(ctx context.Context, input domain.GetStatementInput) (*domain.GetStatementOutput, error)
	AccountGet(ctx context.Context, input domain.AccountTenantIDInput) (*domain.AccountOutput, error)
	AccountCreate(ctx context.Context, input domain.CreateAccountInput) (*domain.AccountOutput, error)
	AccountDeposit(ctx context.Context, input domain.AccountTenantIDWithAmountInput) (*domain.AccountOutput, error)
	AccountWithdraw(ctx context.Context, input domain.AccountTenantIDWithAmountInput) (*domain.AccountOutput, error)
	AccountCheckBalance(ctx context.Context, input domain.AccountTenantIDWithAmountInput) (*domain.AccountCheckOutput, error)
	AssetGetAll(ctx context.Context, input domain.AssetGetAllInput) (*domain.PaginatedAssetOutput, error)
	AssetGetByTicker(ctx context.Context, input domain.AssetGetByTickerInput) (*domain.AssetOutput, error)
	AssetGetPricesByTicker(ctx context.Context, input domain.AssetGetPricesByTickerInput) (*domain.ManyAssetPriceOutput, error)
	AssetCreateAsset(ctx context.Context, input domain.AssetCreateAssetInput) (*domain.AssetOutput, error)
	AssetUpdateAssetPriceByTicker(ctx context.Context, input domain.AssetUpdateAssetPriceByTickerInput) (*domain.AssetOutput, error)
	AssetUpdateAsset(ctx context.Context, input domain.AssetUpdateAssetInput) (*domain.AssetOutput, error)
	AssetDelete(ctx context.Context, input domain.AssetDeleteInput) (*domain.AssetOutput, error)
}
