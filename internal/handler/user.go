package handler

import (
	"net/http"
	"time"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type user struct {
	userSrv IUserService

	appENV             string
	jwtAccessLifetime  int
	jwtRefreshLifetime int
}

func newUser(userSrv IUserService, cfg *config.Config) (*user, error) {
	accessDuration, err := time.ParseDuration(cfg.Server.JWTAccessLifetime)
	if err != nil {
		return nil, err
	}

	refreshDuration, err := time.ParseDuration(cfg.Server.JWTRefreshLifetime)
	if err != nil {
		return nil, err
	}

	return &user{
		userSrv:            userSrv,
		appENV:             cfg.App.ENV,
		jwtAccessLifetime:  int(accessDuration.Seconds()),
		jwtRefreshLifetime: int(refreshDuration.Seconds()),
	}, nil
}

func (u *user) registerRoutes(r chi.Router) {
	r.Post("/users/login", u.login)
	r.Post("/users/registration", u.registration)
	r.Put("/users/refresh", u.refreshTokens)
}

func (u *user) registration(rw http.ResponseWriter, r *http.Request) {
	body, err := parseJSONBody[createUserDTO](r)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}
	defer r.Body.Close()

	access, refresh, err := u.userSrv.Registration(r.Context(), body.Username, body.Password)
	if err != nil {
		if err == apperrors.ErrAlreadyExist {
			Error(rw, HTTPError{Code: http.StatusConflict, Message: "user already exist"})
			return
		}

		logrus.Errorf("registration user err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	u.setCookieTokens(rw, access, refresh)

	rw.WriteHeader(http.StatusCreated)
}

func (u *user) login(rw http.ResponseWriter, r *http.Request) {
	body, err := parseJSONBody[loginUserDTO](r)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}
	defer r.Body.Close()

	access, refresh, err := u.userSrv.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		if err == apperrors.ErrInvalidLoginOrPassword || err == apperrors.ErrNotFound {
			Error(rw, HTTPError{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}

		logrus.Errorf("registration user err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	u.setCookieTokens(rw, access, refresh)

	rw.WriteHeader(http.StatusOK)
}

func (u *user) refreshTokens(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}

	access, refresh, err := u.userSrv.RefreshTokens(ctx, refreshToken.Value)
	if err != nil {
		logrus.Errorf("refresh token err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	u.setCookieTokens(rw, access, refresh)

	rw.WriteHeader(http.StatusOK)
}

func (u *user) setCookieTokens(rw http.ResponseWriter, access, refresh string) {
	secure := false
	if u.appENV == "prod" {
		secure = true
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   u.jwtAccessLifetime,
	})

	http.SetCookie(rw, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   u.jwtRefreshLifetime,
	})
}
