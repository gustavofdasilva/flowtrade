package services

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	repo ports.AccountRepository
}

func NewAccountService(repo ports.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) Get(id uuid.UUID, userID uuid.UUID) (*domain.Account, error) {
	return s.Get(id, userID)
}

func (s *AccountService) Create(userID uuid.UUID, currency domain.Currency) (*domain.Account, error) {
	account := domain.Account{
		UserID:   userID,
		Balance:  decimal.New(0, 1),
		Currency: currency,
	}

	return s.repo.Create(&account)
}

func (s *AccountService) Deposit(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	//TODO: register ledger_entry

	account, err := s.repo.AddBalance(id, userID, amount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) Withdrawal(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	//TODO: register ledger_entry

	account, err := s.repo.SubtractBalance(id, userID, amount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) CheckBalance(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (bool, error) {
	account, err := s.repo.GetByID(id, userID)
	if err != nil {
		return false, err
	}

	return account.Balance.GreaterThanOrEqual(amount), nil
}
