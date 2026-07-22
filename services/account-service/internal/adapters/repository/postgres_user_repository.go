package repository

import (
	"account-service/internal/core/domain"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repo *PostgresUserRepository) GetByEmail(email string) (user *domain.User, err error) {
	query := `
		SELECT id, username, email, password_hash
		FROM public."users"
		WHERE email = $1
		AND deleted_at IS NULL
	`

	row := repo.db.QueryRow(query, email)

	user = &domain.User{}
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("error to scan row: %v", err)
	}

	return user, nil
}

func (repo *PostgresUserRepository) GetByID(id uuid.UUID) (user *domain.User, err error) {
	query := `
		SELECT id, username, email, password_hash
		FROM public."users"
		WHERE id = $1
		AND deleted_at IS NULL
	`

	row := repo.db.QueryRow(query, id)

	user = &domain.User{}
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("error to scan row: %v", err)
	}

	return user, nil
}

func (repo *PostgresUserRepository) CreateUser(user domain.User) (*domain.User, error) {
	query := `
		INSERT INTO public."users"
		(username,email,password_hash)
		VALUES
		($1,$2,$3)
		RETURNING id;
	`

	row := repo.db.QueryRow(query, user.Username, user.Email, user.Password)

	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("error to create user: %v", err)
	}

	return &user, nil
}

func (repo *PostgresUserRepository) DeleteUserByID(id uuid.UUID) (err error) {
	query := `
		UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := repo.db.Exec(query, id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (repo *PostgresUserRepository) UpdateUser(user domain.User) (err error) {
	query := `
		UPDATE users 
		SET 
			username = COALESCE(NULLIF($1, ''), username),
			email = COALESCE(NULLIF($2, ''), email),
			password_hash = COALESCE(NULLIF($3, ''), password_hash)
		WHERE id = $4
		AND deleted_at IS NULL
	`

	result, err := repo.db.Exec(query, user.Username, user.Email, user.Password, user.ID)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (repo *PostgresUserRepository) IsEmailAlreadyInUse(email string, excludedID *uuid.UUID) (exists bool, err error) {
	query := `
		SELECT EXISTS
		(
			SELECT 1
			FROM "users"
			WHERE email = $1
			AND ($2::uuid IS NULL OR id != $2::uuid)
			AND deleted_at IS NULL
		)
	`

	row := repo.db.QueryRow(query, email, excludedID)

	err = row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error to scan row: %v", err)
	}

	return exists, nil
}

func (repo *PostgresUserRepository) IsUsernameAlreadyInUse(username string, excludedID *uuid.UUID) (exists bool, err error) {
	query := `
		SELECT EXISTS
		(
			SELECT 1
			FROM "users"
			WHERE username = $1
			AND ($2::uuid IS NULL OR id != $2::uuid)
			AND deleted_at IS NULL
		)
	`

	row := repo.db.QueryRow(query, username, excludedID)

	err = row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error to scan row: %v", err)
	}

	return exists, nil
}
