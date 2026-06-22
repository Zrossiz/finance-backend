package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	pgUser IUserRepo

	jwtAccessSecret  string
	jwtRefreshSecret string
}

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

type UserClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func newUser(pgUser IUserRepo, jwtAccessSecret string, jwtRefreshSecret string) *user {
	return &user{
		pgUser:           pgUser,
		jwtAccessSecret:  jwtAccessSecret,
		jwtRefreshSecret: jwtRefreshSecret,
	}
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

	access, refresh, err := u.generateTokens(&user)
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

	access, refresh, err := u.generateTokens(user)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (u *user) generateTokens(user *domain.User) (string, string, error) {
	now := time.Now()

	accessClaims := UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	access, err := accessToken.SignedString([]byte(u.jwtAccessSecret))
	if err != nil {
		return "", "", fmt.Errorf("sign access token err: %w", err)
	}

	refreshClaims := UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refresh, err := refreshToken.SignedString([]byte(u.jwtRefreshSecret))
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token err: %w", err)
	}

	return access, refresh, nil
}
