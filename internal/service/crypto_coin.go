package service

import "context"

type cryptoCoin struct {
	pgCryptoCoin ICryptoCoinRepo
}

func newCryptoCoin(pgCryptoCoin ICryptoCoinRepo) *cryptoCoin {
	return &cryptoCoin{
		pgCryptoCoin: pgCryptoCoin,
	}
}

func (c *cryptoCoin) GetAll(ctx context.Context) ([]string, error) {
	return c.pgCryptoCoin.GetAll(ctx)
}
