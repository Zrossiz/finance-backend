package service

import (
	"context"

	"github.com/sirupsen/logrus"
)

type cryptoRate struct {
	pgCryptoPosition ICryptoPositionRepo
	cacheCryptoRates ICryptoRatesCache
	apiCryptoRates   ICryptoRatesAPI
}

func newCryptoRate(
	pgCryptoPosition ICryptoPositionRepo,
	cacheCryptoRates ICryptoRatesCache,
	apiCryptoRates ICryptoRatesAPI,
) *cryptoRate {
	return &cryptoRate{
		pgCryptoPosition: pgCryptoPosition,
		cacheCryptoRates: cacheCryptoRates,
		apiCryptoRates:   apiCryptoRates,
	}
}

func (c *cryptoRate) RefreshCryptoRatesCache() {
	logrus.Info("starting RefreshCryptoRatesCache cron job...")
	ctx := context.Background()

	coinIDs, err := c.pgCryptoPosition.GetUniqueCoinIDs(ctx)
	if err != nil {
		logrus.Errorf("cron get unique coin ids err: %v", err)
		return
	}

	cryptoRates, err := c.apiCryptoRates.GetByIds(ctx, coinIDs)
	if err != nil {
		logrus.Errorf("cron api crypto rates err: %v", err)
		return
	}

	err = c.cacheCryptoRates.Save(ctx, cryptoRates)
	if err != nil {
		logrus.Errorf("cron save cache crypto rates err: %v", err)
		return
	}

	logrus.Infof("finish RefreshCryptoRatesCache cron job. saved crypto rates count: %v", len(cryptoRates))
}
