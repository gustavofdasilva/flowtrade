package repository

import (
	"account-service/internal/core/domain"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) GetByID(id uuid.UUID, userID uuid.UUID) (*domain.Account, error) {
	query := `SELECT id, user_id, balance, currency, created_at, updated_at FROM accounts WHERE id=$1 AND user_id=$2`

	var account domain.Account
	err := r.db.QueryRow(query, id, userID).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (r *PostgresAccountRepository) Create(account domain.Account) (*domain.Account, error) {
	query := `
		INSERT INTO account 
			(user_id, currency, created_at, updated_at)
		VALUES
			$1, $2, NOW(), NOW()
	`

	_, err := r.db.Exec(query, account.UserID, account.Balance)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *PostgresAccountRepository) AddBalance(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	query := `
		UPDATE account 
		SET balance=balance+$1
		WHERE id=$2
		AND user_id=$3
		RETURNING id, user_id, balance, currency, created_at, updated_at
	`

	var account domain.Account
	err := r.db.QueryRow(query, amount, id, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (r *PostgresAccountRepository) SubtractBalance(id uuid.UUID, userID uuid.UUID, amount decimal.Decimal) (*domain.Account, error) {
	query := `
		UPDATE account 
		SET balance=balance-$1
		WHERE id=$2
		AND user_id=$3
		RETURNING id, user_id, balance, currency, created_at, updated_at
	`

	var account domain.Account
	err := r.db.QueryRow(query, amount, id, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}
