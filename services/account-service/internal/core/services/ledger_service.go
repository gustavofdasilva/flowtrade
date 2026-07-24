package services

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"context"

	"github.com/google/uuid"
)

type LedgerService struct {
	repo ports.LedgerRepository
}

func NewLedgerService(repo ports.LedgerRepository) *LedgerService {
	return &LedgerService{
		repo: repo,
	}
}

func (s *LedgerService) GetStatement(ctx context.Context, accountID uuid.UUID, userID uuid.UUID, page, limit int) ([]domain.LedgerEntry, int, error) {
	offset := (page - 1) * limit

	ledgerEntries, count, err := s.repo.GetStatement(ctx, userID, accountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return ledgerEntries, count, err
}
