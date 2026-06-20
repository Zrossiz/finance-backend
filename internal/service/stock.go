package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
)

type stock struct {
	pgStock IStockRepo
}

func newStock(pgStock IStockRepo) *stock {
	return &stock{pgStock: pgStock}
}

func (s *stock) Create(ctx context.Context, payload domain.Stock) error {
	return s.pgStock.Create(ctx, payload)
}

func (s *stock) Delete(ctx context.Context, id uuid.UUID) error {
	return s.pgStock.Delete(ctx, id)
}

func (s *stock) GetAllByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]domain.Stock, error) {
	return s.pgStock.GetAllByUserID(ctx, userID)
}

func (s *stock) GetOne(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	return s.pgStock.GetOne(ctx, id)
}

func (s *stock) Update(ctx context.Context, payload domain.Stock) error {
	return s.pgStock.Update(ctx, payload)
}
