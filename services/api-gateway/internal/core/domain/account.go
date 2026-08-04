package domain

type RegisterInput struct {
	Email    string
	Password string
	Username string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
	User         UserOutput
}

type UserOutput struct {
	ID       string
	Username string
	Email    string
}

type UpdateUserInput struct {
	ID       string
	Username string
	Email    string
	Password string
}

type RefreshTokenInput struct {
	RefreshToken string
}

type LogoutInput struct {
	Token string
}

type UserIDInput struct {
	ID string
}

type GetStatementInput struct {
	UserId    string
	AccountId string
	Page      int64
	Limit     int64
}

type LedgerEntry struct {
	ID           string
	AccountID    string
	Type         string
	Amount       float64
	BalanceAfter float64
	Description  string
	ReferenceID  string
	CreatedAt    string
}

type GetStatementOutput struct {
	LedgerEntries []LedgerEntry
	Page          int64
	Limit         int64
	Total         int64
	TotalPages    int64
}

type AccountTenantIDInput struct {
	UserID    string
	AccountID string
}

type AccountTenantIDWithAmountInput struct {
	UserID    string
	AccountID string
	Amount    float64
}

type CreateAccountInput struct {
	UserID   string
	Currency string
}

type AccountOutput struct {
	ID        string
	UserID    string
	Balance   float64
	Currency  string
	CreatedAt string
	UpdatedAt string
}

type AccountCheckOutput struct {
	Valid bool
}

type AssetGetAllInput struct {
	Page            int64
	Limit           int64
	AssetTypeFilter []string
}

type AssetOutput struct {
	ID             string
	Ticker         string
	Name           string
	Type           string
	Currency       string
	Active         bool
	CreatedAt      string
	UpdatedAt      string
	Price          float64
	PriceUpdatedAt string
}

type PaginatedAssetOutput struct {
	Assets     []AssetOutput
	Page       int64
	Limit      int64
	Total      int64
	TotalPages int64
}

type AssetGetByTickerInput struct {
	Ticker string
}

type AssetGetPricesFilter struct {
	From     string
	To       string
	Interval string
}

type AssetGetPricesByTickerInput struct {
	Ticker string
	Filter AssetGetPricesFilter
}

type AssetPriceOutput struct {
	ID        string
	AssetID   string
	Price     float64
	Source    string
	CreatedAt string
}

type ManyAssetPriceOutput struct {
	AssetPrices []AssetPriceOutput
}

type AssetCreateAssetInput struct {
	Ticker   string
	Name     string
	Type     string
	Currency string
	Active   bool
}

type AssetUpdateAssetPriceByTickerInput struct {
	Ticker string
	Source string
	Price  float64
}

type AssetUpdateAssetInput struct {
	ID       string
	Ticker   string
	Name     string
	Type     string
	Currency string
	Active   bool
}

type AssetDeleteInput struct {
	ID string
}
