package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
)

type bond struct {
	pgBond IBondRepo
}

func newBond(pgBond IBondRepo) *bond {
	return &bond{pgBond: pgBond}
}

func (b *bond) Create(ctx context.Context, payload domain.Bond) error {
	return b.pgBond.Create(ctx, payload)
}

func (b *bond) Delete(ctx context.Context, id uuid.UUID) error {
	return b.pgBond.Delete(ctx, id)
}

func (b *bond) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Bond, error) {
	return b.pgBond.GetAllByUserID(ctx, userID)
}

func (b *bond) GetOne(ctx context.Context, id uuid.UUID) (*domain.Bond, error) {
	return b.pgBond.GetOne(ctx, id)
}
