package dto

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required"`
}

type AccountResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type AmountRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type CheckBalanceResponse struct {
	Valid bool `json:"valid"`
}

type LedgerEntryResponse struct {
	ID           string  `json:"id"`
	AccountID    string  `json:"account_id"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	BalanceAfter float64 `json:"balance_after"`
	Description  string  `json:"description"`
	ReferenceID  string  `json:"reference_id"`
	CreatedAt    string  `json:"created_at"`
}
