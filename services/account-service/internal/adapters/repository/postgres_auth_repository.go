package repository

import (
	"account-service/internal/core/domain"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (repo *PostgresAuthRepository) SaveRefreshToken(userID uuid.UUID, refreshToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := repo.db.Exec(query, userID, refreshToken, expiresAt)
	if err != nil {
		return fmt.Errorf("error to save refresh token: %v", err)
	}

	return nil
}

func (repo *PostgresAuthRepository) GetRefreshToken(refreshToken string) (*domain.RefreshToken, error) {
	query := `
		SELECT user_id, token, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token = $1
	`

	token := &domain.RefreshToken{}
	err := repo.db.QueryRow(query, refreshToken).Scan(&token.UserID, &refreshToken, &token.ExpiresAt, &token.RevokedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("error to scan row: %v", err)
	}

	return token, nil
}

func (repo *PostgresAuthRepository) RevokeRefreshToken(refreshToken string) error {
	query := `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE token = $1
	`

	result, err := repo.db.Exec(query, refreshToken)
	if err != nil {
		return fmt.Errorf("error to revoke refresh token: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return domain.ErrRefreshTokenNotFound
	}

	return nil
}

func (repo *PostgresAuthRepository) RevokeAllRefreshTokens(userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE user_id = $1
	`

	result, err := repo.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error to revoke all refresh tokens: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return domain.ErrRefreshTokenNotFound
	}

	return nil
}
