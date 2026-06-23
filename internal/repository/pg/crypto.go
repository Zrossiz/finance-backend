package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type cryptoPosition struct {
	conn *sql.DB
}

func newCryptoPosition(conn *sql.DB) *cryptoPosition {
	return &cryptoPosition{conn: conn}
}

func (c *cryptoPosition) Create(ctx context.Context, payload domain.CryptoPosition) error {
	_, err := c.conn.ExecContext(
		ctx, createCryptoPositionQuery, payload.ID,
		payload.UserID, payload.Ticker,
		payload.Amount, payload.AvgPriceUSDCents,
	)
	if err != nil {
		return fmt.Errorf("create crypto position db err: %w", err)
	}

	return nil
}

func (c *cryptoPosition) Update(
	ctx context.Context, id uuid.UUID,
	amount decimal.Decimal, avgPriceUsd *int64,
) error {
	res, err := c.conn.ExecContext(
		ctx, updateCryptoPositionQuery,
		id, amount, avgPriceUsd,
	)
	if err != nil {
		return fmt.Errorf("update crypto position db err: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (c *cryptoPosition) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := c.conn.ExecContext(ctx, deleteCryptoPositionQuery, id)
	if err != nil {
		return fmt.Errorf("update crypto position db err: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (c *cryptoPosition) GetAllByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]domain.CryptoPosition, error) {
	rows, err := c.conn.QueryContext(ctx, getUserCryptoPositionsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("get user crypto positions err: %w", err)
	}
	defer rows.Close()

	var positions []domain.CryptoPosition
	for rows.Next() {
		var position domain.CryptoPosition

		err = rows.Scan(
			&position.ID, &position.UserID,
			&position.Ticker, &position.Amount,
			&position.AvgPriceUSDCents,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user crypto position err: %w", err)
		}

		positions = append(positions, position)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user crypto positions err: %w", err)
	}

	return positions, nil
}

func (c *cryptoPosition) GetOneByID(
	ctx context.Context, id uuid.UUID,
) (*domain.CryptoPosition, error) {
	row := c.conn.QueryRowContext(ctx, getCryptoPositionQuery, id)

	var position domain.CryptoPosition
	err := row.Scan(
		&position.ID, &position.UserID,
		&position.Ticker, &position.Amount,
		&position.AvgPriceUSDCents,
	)
	if err != nil {
		return nil, fmt.Errorf("get crypto position err: %w", err)
	}

	return &position, nil
}
