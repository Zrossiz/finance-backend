package handler

import (
	"net/http"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/Zrossiz/finance-backend/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type cryptoPosition struct {
	cryptoPositionSrv ICryptoPositionService
}

func newCryptoPosition(cryptoPositionSrv ICryptoPositionService) *cryptoPosition {
	return &cryptoPosition{cryptoPositionSrv: cryptoPositionSrv}
}

func (c *cryptoPosition) registerRoutes(r chi.Router, accessSecret string) {
	r.Group(func(protected chi.Router) {
		protected.Use(JWTAuth([]byte(accessSecret)))
		r.Get("/cryptos", c.getAllByUserID)
		r.Post("/cryptos", c.create)
		r.Delete("/cryptos/{id}", c.delete)
		r.Patch("/cryptos/{id}", c.update)
	})
}

func (c *cryptoPosition) getAllByUserID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := GetUserClaims(r)

	positions, err := c.cryptoPositionSrv.GetAllByUserID(ctx, claims.UserID)
	if err != nil {
		logrus.Errorf("get user crypto positions err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	total := c.cryptoPositionSrv.CountTotalByPositions(positions)
	res := getUserCryptoPositionsResDTO{
		Total:     total.String(),
		Positions: positions,
	}

	JSON(rw, http.StatusOK, res)
}

func (c *cryptoPosition) create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := helpers.ParseJSONBody[createCryptoPositionDTO](r.Body)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}

	userClaims := GetUserClaims(r)

	amountDecimal, err := decimal.NewFromString(body.Amount)
	if err != nil {
		Error(rw, HTTPError{Code: http.StatusBadRequest, Message: "invalid decimal format"})
		return
	}

	cryptoPosition := domain.CryptoPosition{
		ID:               uuid.New(),
		UserID:           userClaims.UserID,
		Ticker:           body.Ticker,
		Amount:           amountDecimal,
		AvgPriceUSDCents: body.AvgPriceUSDCents,
	}

	err = c.cryptoPositionSrv.Create(ctx, cryptoPosition)
	if err != nil {
		logrus.Errorf("create crypto position err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (c *cryptoPosition) update(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		Error(rw, HTTPError{Code: http.StatusBadRequest, Message: "invalid id format"})
		return
	}

	body, err := helpers.ParseJSONBody[updateCryptoPositionDTO](r.Body)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}
	defer r.Body.Close()

	parsedAmount, err := decimal.NewFromString(body.Amount)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}

	err = c.cryptoPositionSrv.Update(ctx, parsedUUID, parsedAmount, body.AvgPriceUSD)
	if err != nil {
		if err == apperrors.ErrNotFound {
			Error(rw, ErrNotFound)
			return
		}

		logrus.Errorf("registration user err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (c *cryptoPosition) delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		Error(rw, HTTPError{Code: http.StatusBadRequest, Message: "invalid id format"})
		return
	}

	err = c.cryptoPositionSrv.Delete(ctx, parsedUUID)
	if err != nil {
		logrus.Errorf("delete crypto position err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
