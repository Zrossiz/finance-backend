package cache

import (
	"context"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	CryptoRates cryptoRates
}

func New(conn *redis.Client) *Redis {
	return &Redis{
		CryptoRates: *newCryptoRates(conn),
	}
}

func Connect(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Username: cfg.Redis.User,
		Addr:     fmt.Sprintf("%s:%v", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Protocol: 2,
	})

	cmd := rdb.Ping(context.Background())
	if err := cmd.Err(); err != nil {
		return nil, cmd.Err()
	}

	return rdb, nil
}
