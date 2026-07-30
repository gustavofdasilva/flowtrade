package repository

import (
	"account-service/internal/core/domain"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type PostgresLedgerRepository struct {
	db *sql.DB
}

func NewPostgresLedgerRepository(db *sql.DB) *PostgresLedgerRepository {
	return &PostgresLedgerRepository{db: db}
}

func (repo *PostgresLedgerRepository) CreateLedgerEntry(ctx context.Context, entry *domain.LedgerEntry) error {
	query := `
		INSERT INTO account 
			(account_id, type, amount, balance_after, description, reference_id, created_at)
		VALUES
			$1, $2, $3, $4, $5, $6, NOW()
	`

	_, err := repo.db.ExecContext(ctx, query, entry.AccountID, entry.Type, entry.Amount, entry.BalanceAfter, entry.Description, entry.ReferenceID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresLedgerRepository) GetStatement(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, limit, offset int) ([]domain.LedgerEntry, int, error) {
	query := `
		SELECT
			id,
			account_id,
			type,
			amount,
			balance_after,
			description,
			reference_id,
			created_at,
			COUNT(*) OVER()
		FROM
			ledger_entries
		WHERE account_id=$1
		AND (
			SELECT 1 
			FROM accounts
			WHERE id=$1
			AND user_id=$2
		) = 1
		LIMIT $3 OFFSET $4
		ORDER BY created_at DESC;
	`

	rows, err := repo.db.QueryContext(ctx, query, accountID, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var entries []domain.LedgerEntry
	var total int
	for rows.Next() {
		var entry domain.LedgerEntry
		err = rows.Scan(
			&entry.ID,
			&entry.AccountID,
			&entry.Type,
			&entry.Amount,
			&entry.BalanceAfter,
			&entry.Description,
			&entry.ReferenceID,
			&entry.CreatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, err
		}

		entries = append(entries, entry)
	}

	return entries, total, nil
}
