package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/config"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/Zrossiz/finance-backend/internal/helpers"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	pgUser             IUserRepo
	jwtAccessSecret    []byte
	jwtRefreshSecret   []byte
	jwtAccessLifetime  time.Duration
	jwtRefreshLifetime time.Duration
}

func newUser(pgUser IUserRepo, cfg *config.Config) (*user, error) {
	jwtAccessLifetime, err := time.ParseDuration(cfg.Server.JWTAccessLifetime)
	if err != nil {
		return nil, err
	}

	jwtRefreshLifetime, err := time.ParseDuration(cfg.Server.JWTRefreshLifetime)
	if err != nil {
		return nil, err
	}

	return &user{
		pgUser:             pgUser,
		jwtAccessSecret:    []byte(cfg.Server.JWTAccessSecret),
		jwtRefreshSecret:   []byte(cfg.Server.JWTRefreshSecret),
		jwtAccessLifetime:  jwtAccessLifetime,
		jwtRefreshLifetime: jwtRefreshLifetime,
	}, nil
}

func (u *user) Registration(ctx context.Context, username, password string) (string, string, error) {
	existUser, err := u.pgUser.GetByUsername(ctx, username)
	if err != nil && err != apperrors.ErrNotFound {
		return "", "", err
	}

	if existUser != nil {
		return "", "", apperrors.ErrAlreadyExist
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", "", fmt.Errorf("hash password err: %w", err)
	}

	user := domain.User{
		ID:       uuid.New(),
		Username: username,
		Password: string(hash),
	}

	err = u.pgUser.Create(ctx, user)
	if err != nil {
		return "", "", err
	}

	access, refresh, err := helpers.GenerateTokens(
		&user, u.jwtAccessSecret, u.jwtRefreshSecret,
		u.jwtAccessLifetime, u.jwtRefreshLifetime,
	)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (u *user) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := u.pgUser.GetByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", apperrors.ErrInvalidLoginOrPassword
	}

	access, refresh, err := helpers.GenerateTokens(
		user, u.jwtAccessSecret, u.jwtRefreshSecret,
		u.jwtAccessLifetime, u.jwtRefreshLifetime,
	)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
