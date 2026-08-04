// api-gateway/adapters/grpc/client/account_client.go
package grpcclient

import (
	"api-gateway/internal/core/domain"
	"context"

	pb "github.com/gustavofdasilva/flowtrade/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	client pb.AccountServiceClient
}

func NewAccountClient(addr string) *AccountClient {
	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &AccountClient{client: pb.NewAccountServiceClient(conn)}
}

func (c *AccountClient) Register(ctx context.Context, input domain.RegisterInput) error {
	_, err := c.client.Register(ctx, &pb.RegisterRequest{
		Email:    input.Email,
		Password: input.Password,
		Username: input.Username,
	})

	return err
}

func (c *AccountClient) Login(ctx context.Context, input domain.LoginInput) (*domain.AuthOutput, error) {
	token, err := c.client.Login(ctx, &pb.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	return &domain.AuthOutput{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: domain.UserOutput{
			ID:       token.User.Id,
			Username: token.User.Username,
			Email:    token.User.Email,
		},
	}, nil
}

func (c *AccountClient) RefreshToken(ctx context.Context, input domain.RefreshTokenInput) (*domain.AuthOutput, error) {
	token, err := c.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: input.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &domain.AuthOutput{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: domain.UserOutput{
			ID:       token.User.Id,
			Username: token.User.Username,
			Email:    token.User.Email,
		},
	}, nil
}

func (c *AccountClient) Logout(ctx context.Context, input domain.LogoutInput) error {
	//TODO: Remake pb funcs to support LogoutRequest
	_, err := c.client.Logout(ctx, &pb.RefreshTokenRequest{
		RefreshToken: input.Token,
	})

	return err
}

func (c *AccountClient) LogoutAll(ctx context.Context, input domain.UserIDInput) error {
	_, err := c.client.LogoutAll(ctx, &pb.LogoutAllRequest{
		UserId: input.ID,
	})

	return err
}

func (c *AccountClient) UpdateUser(ctx context.Context, input domain.UpdateUserInput) (*domain.UserOutput, error) {
	user, err := c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:       input.ID,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	newUser := domain.UserOutput{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}

	return &newUser, err
}

func (c *AccountClient) GetMe(ctx context.Context, input domain.UserIDInput) (*domain.UserOutput, error) {
	userInfo, err := c.client.GetMe(ctx, &pb.UserIDRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, err
	}

	user := domain.UserOutput{
		ID:       userInfo.Id,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}

	return &user, err
}

func (c *AccountClient) DeleteUser(ctx context.Context, input domain.UserIDInput) error {
	_, err := c.client.DeleteUser(ctx, &pb.UserIDRequest{
		Id: input.ID,
	})

	return err
}

func (c *AccountClient) GetStatement(ctx context.Context, input domain.GetStatementInput) (*domain.GetStatementOutput, error) {
	statement, err := c.client.GetStatement(ctx, &pb.GetStatementRequest{
		UserId:    input.UserId,
		AccountId: input.AccountId,
		Page:      input.Page,
		Limit:     input.Limit,
	})
	if err != nil {
		return nil, err
	}

	ledgerEntries := make([]domain.LedgerEntry, 0, len(statement.LedgerEntries))
	for _, entry := range statement.LedgerEntries {
		ledgerEntries = append(ledgerEntries, domain.LedgerEntry{
			ID:           entry.Id,
			AccountID:    entry.AccountId,
			Type:         entry.Type,
			Amount:       entry.Amount,
			BalanceAfter: entry.BalanceAfter,
			Description:  entry.Description,
			ReferenceID:  entry.ReferenceId,
			CreatedAt:    entry.CreatedAt,
		})
	}

	return &domain.GetStatementOutput{
		LedgerEntries: ledgerEntries,
		Page:          statement.Page,
		Limit:         statement.Limit,
		Total:         statement.Total,
		TotalPages:    statement.TotalPages,
	}, nil
}

func (c *AccountClient) AccountGet(ctx context.Context, input domain.AccountTenantIDInput) (*domain.AccountOutput, error) {
	account, err := c.client.Get(ctx, &pb.AccountTenantIDRequest{
		UserId:    input.UserID,
		AccountId: input.AccountID,
	})
	if err != nil {
		return nil, err
	}

	return toAccountOutput(account), nil
}

func (c *AccountClient) AccountCreate(ctx context.Context, input domain.CreateAccountInput) (*domain.AccountOutput, error) {
	account, err := c.client.Create(ctx, &pb.CreateAccountRequest{
		UserId:   input.UserID,
		Currency: input.Currency,
	})
	if err != nil {
		return nil, err
	}

	return toAccountOutput(account), nil
}

func (c *AccountClient) AccountDeposit(ctx context.Context, input domain.AccountTenantIDWithAmountInput) (*domain.AccountOutput, error) {
	account, err := c.client.Deposit(ctx, &pb.AccountTenantIDWithAmountRequest{
		UserId:    input.UserID,
		AccountId: input.AccountID,
		Amount:    input.Amount,
	})
	if err != nil {
		return nil, err
	}

	return toAccountOutput(account), nil
}

func (c *AccountClient) AccountWithdraw(ctx context.Context, input domain.AccountTenantIDWithAmountInput) (*domain.AccountOutput, error) {
	account, err := c.client.Withdrawal(ctx, &pb.AccountTenantIDWithAmountRequest{
		UserId:    input.UserID,
		AccountId: input.AccountID,
		Amount:    input.Amount,
	})
	if err != nil {
		return nil, err
	}

	return toAccountOutput(account), nil
}

func (c *AccountClient) AccountCheckBalance(ctx context.Context, input domain.AccountTenantIDWithAmountInput) (*domain.AccountCheckOutput, error) {
	check, err := c.client.CheckBalance(ctx, &pb.AccountTenantIDWithAmountRequest{
		UserId:    input.UserID,
		AccountId: input.AccountID,
		Amount:    input.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &domain.AccountCheckOutput{Valid: check.Valid}, nil
}

func (c *AccountClient) AssetGetAll(ctx context.Context, input domain.AssetGetAllInput) (*domain.PaginatedAssetOutput, error) {
	assets, err := c.client.GetAll(ctx, &pb.AssetGetAllRequest{
		Page:            input.Page,
		Limit:           input.Limit,
		AssetTypeFilter: input.AssetTypeFilter,
	})
	if err != nil {
		return nil, err
	}

	assetOutputs := make([]domain.AssetOutput, 0, len(assets.Assets))
	for _, asset := range assets.Assets {
		assetOutputs = append(assetOutputs, toAssetOutput(asset))
	}

	return &domain.PaginatedAssetOutput{
		Assets:     assetOutputs,
		Page:       assets.Page,
		Limit:      assets.Limit,
		Total:      assets.Total,
		TotalPages: assets.TotalPages,
	}, nil
}

func (c *AccountClient) AssetGetByTicker(ctx context.Context, input domain.AssetGetByTickerInput) (*domain.AssetOutput, error) {
	asset, err := c.client.GetByTicker(ctx, &pb.AssetGetByTickerRequest{
		Ticker: input.Ticker,
	})
	if err != nil {
		return nil, err
	}

	output := toAssetOutput(asset)
	return &output, nil
}

func (c *AccountClient) AssetGetPricesByTicker(ctx context.Context, input domain.AssetGetPricesByTickerInput) (*domain.ManyAssetPriceOutput, error) {
	prices, err := c.client.GetPricesByTicker(ctx, &pb.AssetGetPricesByTicker{
		Ticker: input.Ticker,
		Filter: &pb.AssetGetPricesFilter{
			From:     input.Filter.From,
			To:       input.Filter.To,
			Interval: input.Filter.Interval,
		},
	})
	if err != nil {
		return nil, err
	}

	assetPrices := make([]domain.AssetPriceOutput, 0, len(prices.AssetPrices))
	for _, price := range prices.AssetPrices {
		assetPrices = append(assetPrices, domain.AssetPriceOutput{
			ID:        price.Id,
			AssetID:   price.AssetId,
			Price:     price.Price,
			Source:    price.Source,
			CreatedAt: price.CreatedAt,
		})
	}

	return &domain.ManyAssetPriceOutput{AssetPrices: assetPrices}, nil
}

func (c *AccountClient) AssetCreateAsset(ctx context.Context, input domain.AssetCreateAssetInput) (*domain.AssetOutput, error) {
	asset, err := c.client.CreateAsset(ctx, &pb.AssetCreateAssetRequest{
		Ticker:   input.Ticker,
		Name:     input.Name,
		Type:     input.Type,
		Currency: input.Currency,
		Active:   input.Active,
	})
	if err != nil {
		return nil, err
	}

	output := toAssetOutput(asset)
	return &output, nil
}

func (c *AccountClient) AssetUpdateAssetPriceByTicker(ctx context.Context, input domain.AssetUpdateAssetPriceByTickerInput) (*domain.AssetOutput, error) {
	asset, err := c.client.UpdateAssetPriceByTicker(ctx, &pb.AssetUpdateAssetPriceByTickerRequest{
		Ticker: input.Ticker,
		Source: input.Source,
		Price:  input.Price,
	})
	if err != nil {
		return nil, err
	}

	output := toAssetOutput(asset)
	return &output, nil
}

func (c *AccountClient) AssetUpdateAsset(ctx context.Context, input domain.AssetUpdateAssetInput) (*domain.AssetOutput, error) {
	asset, err := c.client.UpdateAsset(ctx, &pb.AssetUpdateAssetRequest{
		Id:       input.ID,
		Ticker:   input.Ticker,
		Name:     input.Name,
		Type:     input.Type,
		Currency: input.Currency,
		Active:   input.Active,
	})
	if err != nil {
		return nil, err
	}

	output := toAssetOutput(asset)
	return &output, nil
}

func (c *AccountClient) AssetDelete(ctx context.Context, input domain.AssetDeleteInput) (*domain.AssetOutput, error) {
	asset, err := c.client.Delete(ctx, &pb.AssetDeleteRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, err
	}

	output := toAssetOutput(asset)
	return &output, nil
}

func toAccountOutput(account *pb.AccountResponse) *domain.AccountOutput {
	return &domain.AccountOutput{
		ID:        account.Id,
		UserID:    account.UserId,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

func toAssetOutput(asset *pb.AssetResponse) domain.AssetOutput {
	return domain.AssetOutput{
		ID:             asset.Id,
		Ticker:         asset.Ticker,
		Name:           asset.Name,
		Type:           asset.Type,
		Currency:       asset.Currency,
		Active:         asset.Active,
		CreatedAt:      asset.CreatedAt,
		UpdatedAt:      asset.UpdatedAt,
		Price:          asset.Price,
		PriceUpdatedAt: asset.PriceUpdatedAt,
	}
}
