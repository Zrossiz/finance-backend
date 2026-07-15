package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

type cryptoRates struct {
	conn *redis.Client
}

const (
	cryptoRatesKey = "crypto:rates"
	cryptoRatesTTL = 1 * time.Minute
)

func newCryptoRates(conn *redis.Client) *cryptoRates {
	return &cryptoRates{conn: conn}
}

func (c *cryptoRates) Save(ctx context.Context, rates domain.CryptoRates) error {
	values := make([]any, 0, len(rates)*2)

	for id, rate := range rates {
		values = append(values, id, rate.USD)
	}

	if len(values) == 0 {
		return nil
	}

	if err := c.conn.HSet(ctx, cryptoRatesKey, values...).Err(); err != nil {
		return err
	}

	return c.conn.Expire(ctx, cryptoRatesKey, cryptoRatesTTL).Err()
}

func (c *cryptoRates) GetByIDs(ctx context.Context, ids []string) (domain.CryptoRates, error) {
	values, err := c.conn.HMGet(ctx, cryptoRatesKey, ids...).Result()
	if err != nil {
		return nil, err
	}

	res := make(domain.CryptoRates)

	for i, value := range values {
		if value == nil {
			continue
		}

		price, err := strconv.ParseFloat(value.(string), 32)
		if err != nil {
			continue
		}

		res[ids[i]] = domain.CryptoRate{
			USD: float32(price),
		}
	}

	return res, nil
}
