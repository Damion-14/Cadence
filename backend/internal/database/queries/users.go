package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/damion-14/cadence/backend/internal/models"
)

type UserQueries struct {
	db *sql.DB
}

func NewUserQueries(db *sql.DB) *UserQueries {
	return &UserQueries{db: db}
}

func (q *UserQueries) CreateUser(ctx context.Context, email, passwordHash, username string) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash, username)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, username, created_at, updated_at
	`

	var user models.User
	err := q.db.QueryRowContext(ctx, query, email, passwordHash, username).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (q *UserQueries) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, username, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := q.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (q *UserQueries) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, username, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := q.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
