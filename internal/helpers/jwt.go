package helpers

import (
	"fmt"
	"time"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

func ValidateJWT(tokenStr string, secret []byte) (*domain.JWTUserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.JWTUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.JWTUserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GenerateTokens(
	user *domain.User,
	accessSecret,
	refreshSecret []byte,
	accessTokenTTL,
	refreshTokenTTL time.Duration,
) (string, string, error) {
	now := time.Now()

	accessClaims := domain.JWTUserClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	access, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return "", "", fmt.Errorf("sign access token err: %w", err)
	}

	refreshClaims := domain.JWTUserClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refresh, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token err: %w", err)
	}

	return access, refresh, nil
}
