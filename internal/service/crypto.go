package service

import (
	"context"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type cryptoPosition struct {
	pgCryptoPositon ICryptoPositionRepo
	apiCryptoRate   ICryptoRatesAPI
}

func newCryptoPosition(pgCryptoPositon ICryptoPositionRepo, apiCryptoRate ICryptoRatesAPI) *cryptoPosition {
	return &cryptoPosition{
		pgCryptoPositon: pgCryptoPositon,
		apiCryptoRate:   apiCryptoRate,
	}
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
	var totalUSD decimal.Decimal

	cryptoPositions, err := c.pgCryptoPositon.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	coinIDs := make([]string, 0, len(cryptoPositions))
	for _, position := range cryptoPositions {
		coinIDs = append(coinIDs, position.CoinID)
	}

	cryptoRates, err := c.apiCryptoRate.GetByIds(ctx, coinIDs)
	if err != nil {
		return nil, err
	}

	for i := range cryptoPositions {
		rate, ok := cryptoRates[cryptoPositions[i].CoinID]
		if !ok {
			continue
		}

		decimalUsd, err := decimal.NewFromString(rate.USD)
		if err != nil {
			return nil, fmt.Errorf("parse usd rate for %s: %w", cryptoPositions[i].CoinID, err)
		}

		cryptoPositions[i].TotalPriceUSD = decimalUsd.Mul(cryptoPositions[i].Amount)
		totalUSD = totalUSD.Add(cryptoPositions[i].TotalPriceUSD)

	}

	return cryptoPositions, nil
}

func (c *cryptoPosition) CountTotalByPositions(positions []domain.CryptoPosition) decimal.Decimal {
	var total decimal.Decimal

	for _, v := range positions {
		total.Add(v.TotalPriceUSD)
	}

	return total
}

func (c *cryptoPosition) GetOneByID(ctx context.Context, id uuid.UUID) (*domain.CryptoPosition, error) {
	return c.pgCryptoPositon.GetOneByID(ctx, id)
}

func (c *cryptoPosition) Update(
	ctx context.Context, id uuid.UUID,
	amount decimal.Decimal, avgPriceUsd *int64,
) error {
	return c.pgCryptoPositon.Update(ctx, id, amount, avgPriceUsd)
}
