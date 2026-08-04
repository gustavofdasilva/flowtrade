package handler

import (
	response "api-gateway/internal/adapters/http"
	"api-gateway/internal/adapters/http/dto"
	"api-gateway/internal/core/domain"
	"api-gateway/internal/core/ports"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountHandler struct {
	accountClient ports.AccountClient
}

func NewAccountHandler(client ports.AccountClient) *AccountHandler {
	return &AccountHandler{accountClient: client}
}

func (h *AccountHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err := h.accountClient.Register(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AccountHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.accountClient.Login(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.TokensResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User: dto.UserResponse{
			ID:       res.User.ID,
			Username: res.User.Username,
			Email:    res.User.Email,
		},
	}

	response.Success(c, http.StatusOK, "Login successful", data)
}

func (h *AccountHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	}

	res, err := h.accountClient.RefreshToken(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.TokensResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User: dto.UserResponse{
			ID:       res.User.ID,
			Username: res.User.Username,
			Email:    res.User.Email,
		},
	}

	response.Success(c, http.StatusOK, "Refresh successful", data)
}

func (h *AccountHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.LogoutInput{
		Token: req.Token,
	}

	err := h.accountClient.Logout(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AccountHandler) LogoutAll(c *gin.Context) {
	userID := c.GetString("userID")

	input := domain.UserIDInput{
		ID: userID,
	}

	err := h.accountClient.LogoutAll(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AccountHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")

	input := domain.UserIDInput{
		ID: userID,
	}

	res, err := h.accountClient.GetMe(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.UserResponse{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
	}

	response.Success(c, http.StatusOK, "User info retrieved", data)
}

func (h *AccountHandler) UpdateUser(c *gin.Context) {
	userID := c.GetString("userID")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.UpdateUserInput{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	log.Println(userID)

	res, err := h.accountClient.UpdateUser(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.UserResponse{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
	}

	response.Success(c, http.StatusOK, "User info updated", data)
}

func (h *AccountHandler) DeleteUser(c *gin.Context) {
	userID := c.GetString("userID")

	input := domain.UserIDInput{
		ID: userID,
	}

	err := h.accountClient.DeleteUser(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AccountHandler) GetStatement(c *gin.Context) {
	userID := c.GetString("userID")
	accountID := c.Param("accountId")
	page, limit := parsePagination(c)

	input := domain.GetStatementInput{
		UserId:    userID,
		AccountId: accountID,
		Page:      page,
		Limit:     limit,
	}

	res, err := h.accountClient.GetStatement(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	entries := make([]dto.LedgerEntryResponse, 0, len(res.LedgerEntries))
	for _, entry := range res.LedgerEntries {
		entries = append(entries, dto.LedgerEntryResponse{
			ID:           entry.ID,
			AccountID:    entry.AccountID,
			Type:         entry.Type,
			Amount:       entry.Amount,
			BalanceAfter: entry.BalanceAfter,
			Description:  entry.Description,
			ReferenceID:  entry.ReferenceID,
			CreatedAt:    entry.CreatedAt,
		})
	}

	response.Paginated(c, http.StatusOK, "Statement retrieved", entries, response.Pagination{
		Page:       int(res.Page),
		Limit:      int(res.Limit),
		Total:      int(res.Total),
		TotalPages: int(res.TotalPages),
	})
}

func (h *AccountHandler) AccountGet(c *gin.Context) {
	userID := c.GetString("userID")
	accountID := c.Param("accountId")

	input := domain.AccountTenantIDInput{
		UserID:    userID,
		AccountID: accountID,
	}

	res, err := h.accountClient.AccountGet(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Account retrieved", toAccountResponse(res))
}

func (h *AccountHandler) AccountCreate(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.CreateAccountInput{
		UserID:   userID,
		Currency: req.Currency,
	}

	res, err := h.accountClient.AccountCreate(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusCreated, "Account created", toAccountResponse(res))
}

func (h *AccountHandler) AccountDeposit(c *gin.Context) {
	userID := c.GetString("userID")
	accountID := c.Param("accountId")

	var req dto.AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.AccountTenantIDWithAmountInput{
		UserID:    userID,
		AccountID: accountID,
		Amount:    req.Amount,
	}

	res, err := h.accountClient.AccountDeposit(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Deposit successful", toAccountResponse(res))
}

func (h *AccountHandler) AccountWithdraw(c *gin.Context) {
	userID := c.GetString("userID")
	accountID := c.Param("accountId")

	var req dto.AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.AccountTenantIDWithAmountInput{
		UserID:    userID,
		AccountID: accountID,
		Amount:    req.Amount,
	}

	res, err := h.accountClient.AccountWithdraw(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Withdrawal successful", toAccountResponse(res))
}

func (h *AccountHandler) AccountCheckBalance(c *gin.Context) {
	userID := c.GetString("userID")
	accountID := c.Param("accountId")

	var req dto.AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.AccountTenantIDWithAmountInput{
		UserID:    userID,
		AccountID: accountID,
		Amount:    req.Amount,
	}

	res, err := h.accountClient.AccountCheckBalance(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Balance checked", dto.CheckBalanceResponse{Valid: res.Valid})
}

func (h *AccountHandler) AssetGetAll(c *gin.Context) {
	page, limit := parsePagination(c)

	input := domain.AssetGetAllInput{
		Page:            page,
		Limit:           limit,
		AssetTypeFilter: c.QueryArray("type"),
	}

	res, err := h.accountClient.AssetGetAll(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	assets := make([]dto.AssetResponse, 0, len(res.Assets))
	for _, asset := range res.Assets {
		assets = append(assets, toAssetResponse(asset))
	}

	response.Paginated(c, http.StatusOK, "Assets retrieved", assets, response.Pagination{
		Page:       int(res.Page),
		Limit:      int(res.Limit),
		Total:      int(res.Total),
		TotalPages: int(res.TotalPages),
	})
}

func (h *AccountHandler) AssetGetByTicker(c *gin.Context) {
	input := domain.AssetGetByTickerInput{
		Ticker: c.Param("ticker"),
	}

	res, err := h.accountClient.AssetGetByTicker(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Asset retrieved", toAssetResponse(*res))
}

func (h *AccountHandler) AssetGetPricesByTicker(c *gin.Context) {
	input := domain.AssetGetPricesByTickerInput{
		Ticker: c.Param("ticker"),
		Filter: domain.AssetGetPricesFilter{
			From:     c.Query("from"),
			To:       c.Query("to"),
			Interval: c.Query("interval"),
		},
	}

	res, err := h.accountClient.AssetGetPricesByTicker(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	prices := make([]dto.AssetPriceResponse, 0, len(res.AssetPrices))
	for _, price := range res.AssetPrices {
		prices = append(prices, dto.AssetPriceResponse{
			ID:        price.ID,
			AssetID:   price.AssetID,
			Price:     price.Price,
			Source:    price.Source,
			CreatedAt: price.CreatedAt,
		})
	}

	response.Success(c, http.StatusOK, "Asset prices retrieved", prices)
}

func (h *AccountHandler) AssetCreateAsset(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.AssetCreateAssetInput{
		Ticker:   req.Ticker,
		Name:     req.Name,
		Type:     req.Type,
		Currency: req.Currency,
		Active:   req.Active,
	}

	res, err := h.accountClient.AssetCreateAsset(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusCreated, "Asset created", toAssetResponse(*res))
}

func (h *AccountHandler) AssetUpdateAssetPriceByTicker(c *gin.Context) {
	var req dto.UpdateAssetPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.AssetUpdateAssetPriceByTickerInput{
		Ticker: c.Param("ticker"),
		Source: req.Source,
		Price:  req.Price,
	}

	res, err := h.accountClient.AssetUpdateAssetPriceByTicker(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Asset price updated", toAssetResponse(*res))
}

func (h *AccountHandler) AssetUpdateAsset(c *gin.Context) {
	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.AssetUpdateAssetInput{
		ID:       c.Param("id"),
		Ticker:   req.Ticker,
		Name:     req.Name,
		Type:     req.Type,
		Currency: req.Currency,
		Active:   req.Active,
	}

	res, err := h.accountClient.AssetUpdateAsset(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Asset updated", toAssetResponse(*res))
}

func (h *AccountHandler) AssetDelete(c *gin.Context) {
	input := domain.AssetDeleteInput{
		ID: c.Param("id"),
	}

	res, err := h.accountClient.AssetDelete(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	response.Success(c, http.StatusOK, "Asset deleted", toAssetResponse(*res))
}

func parsePagination(c *gin.Context) (int64, int64) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	if err != nil || limit < 1 {
		limit = 20
	}

	return page, limit
}

func toAccountResponse(account *domain.AccountOutput) dto.AccountResponse {
	return dto.AccountResponse{
		ID:        account.ID,
		UserID:    account.UserID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

func toAssetResponse(asset domain.AssetOutput) dto.AssetResponse {
	return dto.AssetResponse{
		ID:             asset.ID,
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

func grpcStatusToHTTP(err error) int {
	switch status.Code(err) {
	case codes.NotFound:
		return 404
	case codes.InvalidArgument:
		return 400
	case codes.Unauthenticated:
		return 401
	case codes.PermissionDenied:
		return 403
	case codes.AlreadyExists:
		return 409
	case codes.ResourceExhausted:
		return 429
	default:
		return 500
	}
}
