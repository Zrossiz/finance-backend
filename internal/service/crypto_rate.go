package service

import (
	"context"

	"github.com/sirupsen/logrus"
)

type cryptoRate struct {
	pgCryptoCoin     ICryptoCoinRepo
	cacheCryptoRates ICryptoRatesCache
	apiCryptoRates   ICryptoRatesAPI
}

func newCryptoRate(
	pgCryptoCoin ICryptoCoinRepo,
	apiCryptoRates ICryptoRatesAPI,
	cacheCryptoRates ICryptoRatesCache,
) *cryptoRate {
	return &cryptoRate{
		pgCryptoCoin:     pgCryptoCoin,
		apiCryptoRates:   apiCryptoRates,
		cacheCryptoRates: cacheCryptoRates,
	}
}

func (c *cryptoRate) RefreshCryptoRatesCache() {
	logrus.Info("starting RefreshCryptoRatesCache cron job...")
	ctx := context.Background()

	coins, err := c.pgCryptoCoin.GetAll(ctx)
	if err != nil {
		logrus.Errorf("cron get crypto coin ids err: %v", err)
		return
	}

	coinIDs := make([]string, len(coins))
	for i, v := range coins {
		coinIDs[i] = v.CoinID
	}

	cryptoRates, err := c.apiCryptoRates.GetByIDs(ctx, coinIDs)
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
