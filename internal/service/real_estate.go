package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
)

type realEstate struct {
	pgRealEstate IRealEstateRepo
}

func newRealEstate(pgRealEstate IRealEstateRepo) *realEstate {
	return &realEstate{pgRealEstate: pgRealEstate}
}

func (r *realEstate) Create(ctx context.Context, payload domain.RealEstate) error {
	return r.pgRealEstate.Create(ctx, payload)
}

func (r *realEstate) Delete(ctx context.Context, id uuid.UUID) error {
	return r.pgRealEstate.Delete(ctx, id)
}

func (r *realEstate) GetAllByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]domain.RealEstate, error) {
	return r.pgRealEstate.GetAllByUserID(ctx, userID)
}

func (r *realEstate) GetOne(ctx context.Context, id uuid.UUID) (*domain.RealEstate, error) {
	return r.pgRealEstate.GetOne(ctx, id)
}

func (r *realEstate) Update(ctx context.Context, payload domain.RealEstate) error {
	return r.pgRealEstate.Update(ctx, payload)
}
