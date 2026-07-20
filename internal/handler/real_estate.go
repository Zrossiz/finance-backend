package handler

import (
	"net/http"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/Zrossiz/finance-backend/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type realEstate struct {
	realEstateSrv IRealEstateService
}

func newRealEstate(realEstateSrv IRealEstateService) *realEstate {
	return &realEstate{
		realEstateSrv: realEstateSrv,
	}
}

func (re *realEstate) registerRoutes(r chi.Router, accessSecret string) {
	r.Group(func(protected chi.Router) {
		protected.Use(JWTAuth([]byte(accessSecret)))
		protected.Get("/real-estates", re.getAllByUserID)
		protected.Post("/real-estates", re.create)
		protected.Delete("/real-estates/{id}", re.delete)
		protected.Put("/real-estates/{id}", re.update)
	})
}

func (re *realEstate) create(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := helpers.ParseJSONBody[createRealEstateDTO](r.Body)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}

	userClaims := GetUserClaims(r)

	err = re.realEstateSrv.Create(ctx, domain.RealEstate{
		ID:                 uuid.New(),
		UserID:             userClaims.UserID,
		Name:               body.Name,
		Currency:           body.Currency,
		PurchasePriceCents: body.PurchasePriceCents,
		Purchased:          body.Purchased,
		MonthlyIncomeCents: body.MonthlyIncomeCents,
	})
	if err != nil {
		logrus.Errorf("create real estate err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (re *realEstate) getAllByUserID(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userClaims := GetUserClaims(r)

	realEstates, err := re.realEstateSrv.GetAllByUserID(ctx, userClaims.UserID)
	if err != nil {
		logrus.Errorf("get all real estates by user id err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	JSON(rw, http.StatusOK, realEstates)
}

func (re *realEstate) delete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		Error(rw, HTTPError{Code: http.StatusBadRequest, Message: "invalid id format"})
		return
	}

	err = re.realEstateSrv.Delete(ctx, parsedUUID)
	if err != nil {
		logrus.Errorf("delete real estate err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (re *realEstate) update(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		Error(rw, HTTPError{Code: http.StatusBadRequest, Message: "invalid id format"})
		return
	}

	body, err := helpers.ParseJSONBody[updateRealEstateDTO](r.Body)
	if err != nil {
		Error(rw, ErrBadRequest)
		return
	}
	defer r.Body.Close()

	userClaims := GetUserClaims(r)

	err = re.realEstateSrv.Update(ctx, domain.RealEstate{
		ID:                 parsedUUID,
		UserID:             userClaims.UserID,
		Name:               body.Name,
		Currency:           body.Currency,
		PurchasePriceCents: body.PurchasePriceCents,
		MonthlyIncomeCents: body.MonthlyIncomeCents,
		Purchased:          body.Purchased,
	})
	if err != nil {
		logrus.Errorf("update real estate err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
