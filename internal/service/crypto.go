package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type cryptoPosition struct {
	pgCryptoPositon ICryptoPositionRepo
}

func newCryptoPosition(pgCryptoPositon ICryptoPositionRepo) *cryptoPosition {
	return &cryptoPosition{pgCryptoPositon: pgCryptoPositon}
}

func (c *cryptoPosition) Create(ctx context.Context, payload domain.CryptoPosition) error {
	return c.pgCryptoPositon.Create(ctx, payload)
}

func (c *cryptoPosition) Delete(ctx context.Context, id uuid.UUID) error {
	return c.pgCryptoPositon.Delete(ctx, id)
}

func (c *cryptoPosition) GetAllByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]domain.CryptoPosition, error) {
	return c.pgCryptoPositon.GetAllByUserID(ctx, userID)
}

func (c *cryptoPosition) GetOneByID(ctx context.Context, id uuid.UUID) (*domain.CryptoPosition, error) {
	return c.pgCryptoPositon.GetOneByID(ctx, id)
}

func (c *cryptoPosition) Update(
	ctx context.Context, id uuid.UUID,
	amount decimal.Decimal, avgPriceUsd int64,
) error {
	return c.pgCryptoPositon.Update(ctx, id, amount, avgPriceUsd)
}
