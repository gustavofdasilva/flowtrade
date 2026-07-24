package database

import (
	"context"
	"database/sql"
)

type PostgresTransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *PostgresTransactionManager {
	return &PostgresTransactionManager{db: db}
}

func (m *PostgresTransactionManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
