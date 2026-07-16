package pgrepo

import (
	"context"
	"database/sql"
	"fmt"
)

type cryptoCoin struct {
	conn *sql.DB
}

func newCryptoCoin(conn *sql.DB) *cryptoCoin {
	return &cryptoCoin{conn: conn}
}

func (c *cryptoCoin) GetAll(ctx context.Context) ([]string, error) {
	rows, err := c.conn.QueryContext(ctx, getCoinsQuery)
	if err != nil {
		return nil, fmt.Errorf("get crypto coins db err: %w", err)
	}

	var coins []string
	for rows.Next() {
		var coin string
		err = rows.Scan(&coin)
		if err != nil {
			return nil, fmt.Errorf("scan crypto coins err: %w", err)
		}

		coins = append(coins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("crypto coins iterate err: %w", err)
	}

	return coins, nil
}
