package middleware

import (
	"context"
	"net/http"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/Zrossiz/finance-backend/internal/handler"
	"github.com/Zrossiz/finance-backend/internal/helpers"
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

func JWTAuth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				handler.Error(w, handler.ErrUnauthorized)
				return
			}

			claims, err := helpers.ValidateJWT(cookie.Value, secret)
			if err != nil {
				handler.Error(w, handler.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserClaims(r *http.Request) *domain.JWTUserClaims {
	if claims, ok := r.Context().Value(UserClaimsKey).(*domain.JWTUserClaims); ok {
		return claims
	}
	return nil
}
