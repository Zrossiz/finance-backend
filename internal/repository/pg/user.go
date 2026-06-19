package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
)

type user struct {
	conn *sql.DB
}

func newUser(conn *sql.DB) *user {
	return &user{
		conn: conn,
	}
}

func (u *user) Create(ctx context.Context, payload domain.User) error {
	_, err := u.conn.ExecContext(
		ctx, createUserQuery,
		payload.ID, payload.Username,
		payload.Password,
	)
	if err != nil {
		return fmt.Errorf("create user db err: %w", err)
	}

	return nil
}

func (u *user) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := u.conn.QueryRowContext(ctx, getUserByUsernameQuery, username)

	var user domain.User
	err := row.Scan(
		&user.ID, &user.Username,
		&user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("create user db err: %w", err)
	}

	return &user, nil
}
