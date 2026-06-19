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
	return nil, nil
}

func (s *stock) GetOne(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	_, err := s.pgStock.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *stock) Update(ctx context.Context, payload domain.Stock) error {
	return nil
}
