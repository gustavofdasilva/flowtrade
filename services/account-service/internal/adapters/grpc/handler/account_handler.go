package handler

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"context"
	"time"

	"github.com/google/uuid"
	pb "github.com/gustavofdasilva/flowtrade/proto/pb"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountGRPCHandler struct {
	pb.UnimplementedAccountServiceServer
	authSvc    ports.AuthService
	userSvc    ports.UserService
	ledgerSvc  ports.LedgerService
	accountSvc ports.AccountService
	assetSvc   ports.AssetService
}

func NewAccountGRPCHandler(authSvc ports.AuthService, userSvc ports.UserService, ledgerSvc ports.LedgerService, accountSvc ports.AccountService, assetSvc ports.AssetService) *AccountGRPCHandler {
	return &AccountGRPCHandler{
		authSvc:    authSvc,
		userSvc:    userSvc,
		ledgerSvc:  ledgerSvc,
		accountSvc: accountSvc,
		assetSvc:   assetSvc,
	}
}

func (h *AccountGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	token, err := h.authSvc.Login(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	user := pb.UserInfo{
		Id:       token.User.ID.String(),
		Username: token.User.Username,
		Email:    token.User.Email,
	}

	return &pb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         &user,
	}, nil
}

func (h *AccountGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {

	token, err := h.authSvc.Refresh(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	user := pb.UserInfo{
		Id:       token.User.ID.String(),
		Username: token.User.Username,
		Email:    token.User.Email,
	}

	return &pb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         &user,
	}, nil
}

func (h *AccountGRPCHandler) Logout(ctx context.Context, req *pb.RefreshTokenRequest) (*emptypb.Empty, error) {

	err := h.authSvc.Logout(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *AccountGRPCHandler) LogoutAll(ctx context.Context, req *pb.LogoutAllRequest) (*emptypb.Empty, error) {

	err := h.authSvc.LogoutAll(uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *AccountGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	user := domain.User{
		Username: req.Username,
		Email:    req.Username,
		Password: req.Password,
	}

	_, err := h.userSvc.Register(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *AccountGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserInfo, error) {
	user := domain.User{
		ID:       uuid.MustParse(req.Id),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	newUser, err := h.userSvc.Update(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.UserInfo{
		Id:       user.ID.String(),
		Username: newUser.Username,
		Email:    newUser.Email,
	}, nil
}

func (h *AccountGRPCHandler) GetMe(ctx context.Context, req *pb.UserIDRequest) (*pb.UserInfo, error) {
	user, err := h.userSvc.GetByID(uuid.MustParse(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.UserInfo{
		Id:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *AccountGRPCHandler) DeleteUser(ctx context.Context, req *pb.UserIDRequest) (*emptypb.Empty, error) {
	err := h.userSvc.Delete(uuid.MustParse(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *AccountGRPCHandler) LedgerEntryGetStatement(ctx context.Context, req *pb.GetStatementRequest) (*pb.GetStatementResponse, error) {
	statement, total, err := h.ledgerSvc.GetStatement(ctx, uuid.MustParse(req.UserId), uuid.MustParse(req.AccountId), int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	ledgerEntries := make([]*pb.LedgerEntry, len(statement))
	for i, entry := range statement {
		ledgerEntries[i] = &pb.LedgerEntry{
			Id:           entry.ID.String(),
			AccountId:    entry.AccountID.String(),
			Type:         entry.Type.String(),
			Amount:       entry.Amount.InexactFloat64(),
			BalanceAfter: entry.BalanceAfter.InexactFloat64(),
			Description:  entry.Description,
			ReferenceId:  entry.ID.String(),
			CreatedAt:    entry.CreatedAt.String(),
		}
	}

	totalPages := (total + int(req.Limit) - 1) / int(req.Limit)
	return &pb.GetStatementResponse{
		LedgerEntries: ledgerEntries,
		Page:          int64(req.Page),
		Limit:         int64(req.Limit),
		Total:         int64(total),
		TotalPages:    int64(totalPages),
	}, nil
}

func (h *AccountGRPCHandler) AccountGet(ctx context.Context, req *pb.AccountTenantIDRequest) (*pb.AccountResponse, error) {
	account, err := h.accountSvc.Get(ctx, uuid.MustParse(req.AccountId), uuid.MustParse(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Id:        account.ID.String(),
		UserId:    account.UserID.String(),
		Balance:   account.Balance.InexactFloat64(),
		Currency:  string(account.Currency),
		CreatedAt: account.CreatedAt.String(),
		UpdatedAt: account.UpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AccountCreate(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	account, err := h.accountSvc.Create(ctx, uuid.MustParse(req.UserId), domain.Currency(req.Currency))
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Id:        account.ID.String(),
		UserId:    account.UserID.String(),
		Balance:   account.Balance.InexactFloat64(),
		Currency:  string(account.Currency),
		CreatedAt: account.CreatedAt.String(),
		UpdatedAt: account.UpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AccountDeposit(ctx context.Context, req *pb.AccountTenantIDWithAmountRequest) (*pb.AccountResponse, error) {
	account, err := h.accountSvc.Deposit(ctx, uuid.MustParse(req.AccountId), uuid.MustParse(req.UserId), decimal.NewFromFloat(req.Amount))
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Id:        account.ID.String(),
		UserId:    account.UserID.String(),
		Balance:   account.Balance.InexactFloat64(),
		Currency:  string(account.Currency),
		CreatedAt: account.CreatedAt.String(),
		UpdatedAt: account.UpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AccountWithDrawal(ctx context.Context, req *pb.AccountTenantIDWithAmountRequest) (*pb.AccountResponse, error) {
	account, err := h.accountSvc.Withdrawal(ctx, uuid.MustParse(req.AccountId), uuid.MustParse(req.UserId), decimal.NewFromFloat(req.Amount))
	if err != nil {
		return nil, err
	}

	return &pb.AccountResponse{
		Id:        account.ID.String(),
		UserId:    account.UserID.String(),
		Balance:   account.Balance.InexactFloat64(),
		Currency:  string(account.Currency),
		CreatedAt: account.CreatedAt.String(),
		UpdatedAt: account.UpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AccountCheckBalance(ctx context.Context, req *pb.AccountTenantIDWithAmountRequest) (*pb.AccountCheckResponse, error) {
	isValid, err := h.accountSvc.CheckBalance(ctx, uuid.MustParse(req.AccountId), uuid.MustParse(req.UserId), decimal.NewFromFloat(req.Amount))
	if err != nil {
		return nil, err
	}

	return &pb.AccountCheckResponse{
		Valid: isValid,
	}, nil
}

func (h *AccountGRPCHandler) AssetGetAll(ctx context.Context, req *pb.AssetGetAllRequest) (*pb.PaginatedAssetResponse, error) {
	var filtersToType []domain.AssetType

	for _, f := range req.AssetTypeFilter {
		filtersToType = append(filtersToType, domain.AssetType(f))
	}

	assets, total, err := h.assetSvc.GetAll(ctx, int(req.Page), int(req.Limit), filtersToType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	assetsPb := make([]*pb.AssetResponse, len(assets))
	for i, entry := range assets {
		assetsPb[i] = &pb.AssetResponse{
			Id:             entry.ID.String(),
			Ticker:         entry.Ticker,
			Name:           entry.Name,
			Type:           string(entry.Type),
			Currency:       string(entry.Currency),
			Active:         entry.Active,
			CreatedAt:      entry.CreatedAt.String(),
			UpdatedAt:      entry.UpdatedAt.String(),
			Price:          entry.Price.InexactFloat64(),
			PriceUpdatedAt: entry.PriceUpdatedAt.String(),
		}
	}

	totalPages := (total + int(req.Limit) - 1) / int(req.Limit)
	return &pb.PaginatedAssetResponse{
		Assets:     assetsPb,
		Page:       int64(req.Page),
		Limit:      int64(req.Limit),
		Total:      int64(total),
		TotalPages: int64(totalPages),
	}, nil
}

func (h *AccountGRPCHandler) AssetGetByTicker(ctx context.Context, req *pb.AssetGetByTickerRequest) (*pb.AssetResponse, error) {
	asset, err := h.assetSvc.GetByTicker(ctx, req.Ticker)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.AssetResponse{
		Id:             asset.ID.String(),
		Ticker:         asset.Ticker,
		Name:           asset.Name,
		Type:           string(asset.Type),
		Currency:       string(asset.Currency),
		Active:         asset.Active,
		CreatedAt:      asset.CreatedAt.String(),
		UpdatedAt:      asset.UpdatedAt.String(),
		Price:          asset.Price.InexactFloat64(),
		PriceUpdatedAt: asset.PriceUpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AssetGetPricesByTicker(ctx context.Context, req *pb.AssetGetPricesByTicker) (*pb.ManyAssetPriceResponse, error) {

	layout := `2006-01-02`
	var from time.Time
	var to time.Time
	var interval time.Duration
	var err error

	filter := domain.AssetPriceFilter{}

	if req.Filter.From != "" {
		var err error
		from, err = time.Parse(layout, req.Filter.From)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		filter.From = &from
	}

	if req.Filter.To != "" {
		to, err = time.Parse(layout, req.Filter.To)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		filter.To = &to
	}

	if req.Filter.Interval != "" {
		interval, err = time.ParseDuration(req.Filter.Interval)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		filter.Interval = &interval
	}

	assetPrices, err := h.assetSvc.GetPricesByTicker(ctx, req.Ticker, &filter)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	assetPricesPb := make([]*pb.AssetPriceResponse, len(assetPrices))
	for i, entry := range assetPrices {
		assetPricesPb[i] = &pb.AssetPriceResponse{
			Id:        entry.ID.String(),
			AssetId:   entry.AssetID.String(),
			CreatedAt: entry.CreatedAt.String(),
			Price:     entry.Price.InexactFloat64(),
			Source:    entry.Source,
		}
	}

	return &pb.ManyAssetPriceResponse{
		AssetPrices: assetPricesPb,
	}, nil
}

func (h *AccountGRPCHandler) AssetCreateAsset(ctx context.Context, req *pb.AssetCreateAssetRequest) (*pb.AssetResponse, error) {

	asset := domain.Asset{
		Ticker:   req.Ticker,
		Name:     req.Name,
		Currency: domain.Currency(req.Currency),
		Type:     domain.AssetType(req.Type),
	}

	newAsset, err := h.assetSvc.CreateAsset(ctx, asset)
	if err != nil {
		return nil, err
	}

	return &pb.AssetResponse{
		Id:             newAsset.ID.String(),
		Ticker:         newAsset.Ticker,
		Name:           newAsset.Name,
		Type:           string(newAsset.Type),
		Currency:       string(newAsset.Currency),
		Active:         newAsset.Active,
		CreatedAt:      newAsset.CreatedAt.String(),
		UpdatedAt:      newAsset.UpdatedAt.String(),
		Price:          newAsset.Price.InexactFloat64(),
		PriceUpdatedAt: newAsset.PriceUpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AssetUpdateAssetPriceByTicker(ctx context.Context, req *pb.AssetUpdateAssetPriceByTickerRequest) (*pb.AssetResponse, error) {

	newAsset, err := h.assetSvc.UpdateAssetPriceByTicker(ctx, req.Source, req.Ticker, decimal.NewFromFloat(req.Price))
	if err != nil {
		return nil, err
	}

	return &pb.AssetResponse{
		Id:             newAsset.ID.String(),
		Ticker:         newAsset.Ticker,
		Name:           newAsset.Name,
		Type:           string(newAsset.Type),
		Currency:       string(newAsset.Currency),
		Active:         newAsset.Active,
		CreatedAt:      newAsset.CreatedAt.String(),
		UpdatedAt:      newAsset.UpdatedAt.String(),
		Price:          newAsset.Price.InexactFloat64(),
		PriceUpdatedAt: newAsset.PriceUpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AssetUpdateAsset(ctx context.Context, req *pb.AssetUpdateAssetRequest) (*pb.AssetResponse, error) {

	asset := domain.Asset{
		ID:       uuid.MustParse(req.Id),
		Ticker:   req.Ticker,
		Name:     req.Name,
		Currency: domain.Currency(req.Currency),
		Type:     domain.AssetType(req.Type),
		Active:   req.Active,
	}

	newAsset, err := h.assetSvc.UpdateAsset(ctx, asset)
	if err != nil {
		return nil, err
	}

	return &pb.AssetResponse{
		Id:             newAsset.ID.String(),
		Ticker:         newAsset.Ticker,
		Name:           newAsset.Name,
		Type:           string(newAsset.Type),
		Currency:       string(newAsset.Currency),
		Active:         newAsset.Active,
		CreatedAt:      newAsset.CreatedAt.String(),
		UpdatedAt:      newAsset.UpdatedAt.String(),
		Price:          newAsset.Price.InexactFloat64(),
		PriceUpdatedAt: newAsset.PriceUpdatedAt.String(),
	}, nil
}

func (h *AccountGRPCHandler) AssetDelete(ctx context.Context, req *pb.AssetDeleteRequest) (*pb.AssetResponse, error) {

	newAsset, err := h.assetSvc.Delete(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.AssetResponse{
		Id:             newAsset.ID.String(),
		Ticker:         newAsset.Ticker,
		Name:           newAsset.Name,
		Type:           string(newAsset.Type),
		Currency:       string(newAsset.Currency),
		Active:         newAsset.Active,
		CreatedAt:      newAsset.CreatedAt.String(),
		UpdatedAt:      newAsset.UpdatedAt.String(),
		Price:          newAsset.Price.InexactFloat64(),
		PriceUpdatedAt: newAsset.PriceUpdatedAt.String(),
	}, nil
}
