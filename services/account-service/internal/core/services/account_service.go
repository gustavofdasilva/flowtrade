package services

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	txManager   ports.TransactionManager
	accountRepo ports.AccountRepository
	ledgerRepo  ports.LedgerRepository
}

func NewAccountService(txManager ports.TransactionManager, accountRepo ports.AccountRepository, ledgerRepo ports.LedgerRepository) *AccountService {
	return &AccountService{
		txManager:   txManager,
		accountRepo: accountRepo,
		ledgerRepo:  ledgerRepo,
	}
}

func (s *AccountService) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Account, error) {
	return s.accountRepo.GetByID(ctx, id, userID)
}

func (s *AccountService) Create(ctx context.Context, userID uuid.UUID, currency domain.Currency) (*domain.Account, error) {
	account := domain.Account{
		UserID:   userID,
		Balance:  decimal.New(0, 1),
		Currency: currency,
	}

	return s.accountRepo.Create(ctx, &account)
}

func (s *AccountService) Deposit(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {

	var result *domain.Account
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		account, err := s.accountRepo.AddBalance(txCtx, id, userID, amount)
		if err != nil {
			return err
		}

		if err := s.ledgerRepo.CreateLedgerEntry(txCtx, &domain.LedgerEntry{
			AccountID:    account.ID,
			Type:         domain.LedgerTypeDeposit,
			Amount:       amount,
			BalanceAfter: account.Balance,
		}); err != nil {
			return err
		}

		result = account
		return nil
	})

	return result, err
}

func (s *AccountService) Withdrawal(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {

	var result *domain.Account
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		account, err := s.accountRepo.SubtractBalance(txCtx, id, userID, amount)
		if err != nil {
			return err
		}

		if err := s.ledgerRepo.CreateLedgerEntry(txCtx, &domain.LedgerEntry{
			AccountID:    account.ID,
			Type:         domain.LedgerTypeWithdrawal,
			Amount:       amount,
			BalanceAfter: account.Balance,
		}); err != nil {
			return err
		}

		result = account
		return nil
	})

	return result, err
}

func (s *AccountService) CheckBalance(ctx context.Context, id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (bool, error) {
	account, err := s.accountRepo.GetByID(ctx, id, userID)
	if err != nil {
		return false, err
	}

	return account.Balance.GreaterThanOrEqual(amount), nil
}
