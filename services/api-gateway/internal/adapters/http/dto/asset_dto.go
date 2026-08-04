package dto

type AssetResponse struct {
	ID             string  `json:"id"`
	Ticker         string  `json:"ticker"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Currency       string  `json:"currency"`
	Active         bool    `json:"active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	Price          float64 `json:"price"`
	PriceUpdatedAt string  `json:"price_updated_at"`
}

type CreateAssetRequest struct {
	Ticker   string `json:"ticker" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Active   bool   `json:"active"`
}

type UpdateAssetRequest struct {
	Ticker   string `json:"ticker,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Currency string `json:"currency,omitempty"`
	Active   bool   `json:"active"`
}

type UpdateAssetPriceRequest struct {
	Source string  `json:"source" binding:"required"`
	Price  float64 `json:"price" binding:"required,gt=0"`
}

type AssetPriceResponse struct {
	ID        string  `json:"id"`
	AssetID   string  `json:"asset_id"`
	Price     float64 `json:"price"`
	Source    string  `json:"source"`
	CreatedAt string  `json:"created_at"`
}
