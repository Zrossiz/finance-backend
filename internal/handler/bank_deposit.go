package handler

import (
	"net/http"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/Zrossiz/finance-backend/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type bankDeposit struct {
	bankDepositSrv IBankDepositService
}

func newBankDeposit(bankDepositSrv IBankDepositService) *bankDeposit {
	return &bankDeposit{bankDepositSrv: bankDepositSrv}
}

func (b *bankDeposit) registerRoutes(r chi.Router, accessSecret string) {
	r.Group(func(protected chi.Router) {
		protected.Use(JWTAuth([]byte(accessSecret)))
		protected.Get("/bank-deposit", b.getAllByUserID)
		protected.Post("/bank-deposit", b.create)
		protected.Delete("/bank-deposit", b.delete)
	})
}

func (b *bankDeposit) delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		Error(rw, HTTPError{Code: http.StatusBadRequest, Message: "invalid id format"})
		return
	}

	err = b.bankDepositSrv.Delete(ctx, parsedUUID)
	if err != nil {
		logrus.Errorf("delete bank deposit err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (b *bankDeposit) getAllByUserID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userClaims := GetUserClaims(r)

	deposits, err := b.bankDepositSrv.GetAllByUserID(ctx, userClaims.UserID)
	if err != nil {
		logrus.Errorf("get all bank deposits by user id err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	JSON(rw, http.StatusOK, deposits)
}

func (b *bankDeposit) create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := helpers.ParseJSONBody[createBankDepositDTO](r.Body)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}

	userClaims := GetUserClaims(r)

	interestRateDecimal, err := decimal.NewFromString(body.InterestRate)
	if err != nil {
		Error(rw, HTTPError{Code: http.StatusBadRequest, Message: "invalid decimal format"})
		return
	}

	bankDeposit := domain.BankDeposit{
		ID:           uuid.New(),
		UserID:       userClaims.UserID,
		Name:         body.Name,
		Currency:     body.Currency,
		AmountCents:  body.AmountCents,
		OpenedAt:     body.OpenedAt,
		PeriodMonths: body.PeriodMonths,
		InterestRate: interestRateDecimal,
	}

	err = b.bankDepositSrv.Create(ctx, bankDeposit)
	if err != nil {
		logrus.Errorf("create bank deposit err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
