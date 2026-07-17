package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
)

type cryptoCoin struct {
	pgCryptoCoin ICryptoCoinRepo
}

func newCryptoCoin(pgCryptoCoin ICryptoCoinRepo) *cryptoCoin {
	return &cryptoCoin{
		pgCryptoCoin: pgCryptoCoin,
	}
}

func (c *cryptoCoin) GetAll(ctx context.Context) ([]domain.CryptoCoin, error) {
	return c.pgCryptoCoin.GetAll(ctx)
}
