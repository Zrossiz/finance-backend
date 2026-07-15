package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type cryptoPosition struct {
	pgCryptoPositon ICryptoPositionRepo
	apiCryptoRate   ICryptoRatesAPI
	cacheCryptoRate ICryptoRatesCache
}

func newCryptoPosition(
	pgCryptoPositon ICryptoPositionRepo,
	apiCryptoRate ICryptoRatesAPI,
	cacheCryptoRate ICryptoRatesCache,
) *cryptoPosition {
	return &cryptoPosition{
		pgCryptoPositon: pgCryptoPositon,
		apiCryptoRate:   apiCryptoRate,
		cacheCryptoRate: cacheCryptoRate,
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
	cryptoPositions, err := c.pgCryptoPositon.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	coinIDs := make([]string, 0, len(cryptoPositions))
	for _, position := range cryptoPositions {
		if position.Ticker != "usdt" && position.Ticker != "usdc" {
			coinIDs = append(coinIDs, position.CoinID)
		}
	}

	cryptoRates, err := c.cacheCryptoRate.GetByIDs(ctx, coinIDs)
	if err != nil {
		logrus.Errorf("get cache crypto rates err: %v", err)

		cryptoRates, err = c.apiCryptoRate.GetByIds(ctx, coinIDs)
		if err != nil {
			return nil, err
		}

		err = c.cacheCryptoRate.Save(ctx, cryptoRates)
		if err != nil {
			logrus.Errorf("set cache crypto rates err: %v", err)
		}
	}

	if len(cryptoRates) != len(coinIDs) {
		cryptoRates, err = c.apiCryptoRate.GetByIds(ctx, coinIDs)
		if err != nil {
			return nil, err
		}
	}

	for i := range cryptoPositions {
		position := &cryptoPositions[i]

		if position.Ticker == "usdt" || position.Ticker == "usdc" {
			position.TotalPriceUSD = position.Amount
		} else {
			rate, ok := cryptoRates[position.CoinID]
			if !ok {
				continue
			}

			currentPriceUSD := decimal.NewFromFloat32(rate.USD)
			position.TotalPriceUSD = currentPriceUSD.Mul(position.Amount)
		}

		if position.AvgPriceUSDCents != nil {
			avgPriceUSD := decimal.NewFromInt(*position.AvgPriceUSDCents).Div(decimal.NewFromInt(100))
			investedUSD := avgPriceUSD.Mul(position.Amount)

			position.ProfitUSD = position.TotalPriceUSD.Sub(investedUSD)
		}
	}

	return cryptoPositions, nil
}

func (c *cryptoPosition) CountTotalByPositions(positions []domain.CryptoPosition) decimal.Decimal {
	var totalUSD decimal.Decimal

	for i := range positions {
		totalUSD = totalUSD.Add(positions[i].TotalPriceUSD)
	}

	return totalUSD
}

func (c *cryptoPosition) CountTotalProfitByPositions(positions []domain.CryptoPosition) decimal.Decimal {
	var totalUSD decimal.Decimal

	for i := range positions {
		totalUSD = totalUSD.Add(positions[i].ProfitUSD)
	}

	return totalUSD
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
