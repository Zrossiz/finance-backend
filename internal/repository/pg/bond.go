package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
)

type bond struct {
	conn *sql.DB
}

func newBond(conn *sql.DB) *bond {
	return &bond{conn: conn}
}

func (b *bond) Create(ctx context.Context, payload domain.Bond) error {
	_, err := b.conn.ExecContext(
		ctx, createBondQuery,
		payload.ID, payload.UserID,
		payload.Ticker, payload.Currency,
		payload.Amount, payload.AvgPriceCents,
		payload.CouponCents, payload.CouponPeriodMonths,
	)
	if err != nil {
		return fmt.Errorf("create bond db err: %w", err)
	}

	return nil
}

func (b *bond) GetOne(ctx context.Context, id uuid.UUID) (*domain.Bond, error) {
	row := b.conn.QueryRowContext(ctx, getBondQuery, id)

	var bond domain.Bond
	err := row.Scan(
		&bond.ID, &bond.UserID,
		&bond.Ticker, &bond.Currency,
		&bond.Amount, &bond.AvgPriceCents,
		&bond.CouponCents, &bond.CouponPeriodMonths,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("get bond db err: %w", err)
	}

	return &bond, nil
}

func (b *bond) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Bond, error) {
	rows, err := b.conn.QueryContext(ctx, getUserBondsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("get user bonds db err: %w", err)
	}

	var bonds []domain.Bond
	for rows.Next() {
		var bond domain.Bond
		err = rows.Scan(
			&bond.ID, &bond.UserID,
			&bond.Ticker, &bond.Currency,
			&bond.Amount, &bond.AvgPriceCents,
			&bond.CouponCents, &bond.CouponPeriodMonths,
		)
		if err != nil {
			return nil, fmt.Errorf("scan bond err: %w", err)
		}

		bonds = append(bonds, bond)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate user bonds err: %w", err)
	}

	return bonds, nil
}

func (b *bond) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := b.conn.ExecContext(ctx, deleteBondQuery, id)
	if err != nil {
		return fmt.Errorf("delete bond db err: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
