package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
)

type realEstate struct {
	conn *sql.DB
}

func newRealEstate(conn *sql.DB) *realEstate {
	return &realEstate{conn: conn}
}

func (r *realEstate) Create(ctx context.Context, payload domain.RealEstate) error {
	_, err := r.conn.ExecContext(
		ctx, createRealEstateQuery,
		payload.ID, payload.UserID, payload.Name,
		payload.Currency, payload.PurchasePriceCents,
		payload.MonthlyIncomeCents, payload.Purchased,
	)
	if err != nil {
		return fmt.Errorf("create real estate db err: %w", err)
	}

	return nil
}

func (r *realEstate) GetAllByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]domain.RealEstate, error) {
	rows, err := r.conn.QueryContext(ctx, getUserRealEstatesQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("get real estates db err: %w", err)
	}

	var realEstates []domain.RealEstate
	for rows.Next() {
		var realEstate domain.RealEstate

		err := rows.Scan(
			&realEstate.ID, &realEstate.UserID,
			&realEstate.Name, &realEstate.Currency,
			&realEstate.PurchasePriceCents, &realEstate.MonthlyIncomeCents,
			&realEstate.Purchased,
		)
		if err != nil {
			return nil, fmt.Errorf("scan real estate err: %w", err)
		}

		realEstates = append(realEstates, realEstate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user real estates err: %w", err)
	}

	return realEstates, nil
}

func (r *realEstate) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.conn.ExecContext(ctx, deleteRealEstateQuery, id)
	if err != nil {
		return fmt.Errorf("delete real estate db err: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *realEstate) Update(ctx context.Context, payload domain.RealEstate) error {
	res, err := r.conn.ExecContext(
		ctx, updateRealEstateQuery, payload.ID,
		payload.Name, payload.Currency,
		payload.PurchasePriceCents,
		payload.MonthlyIncomeCents, payload.Purchased,
	)
	if err != nil {
		return fmt.Errorf("update real estate err: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *realEstate) GetOne(ctx context.Context, id uuid.UUID) (*domain.RealEstate, error) {
	row := r.conn.QueryRowContext(ctx, getRealEstateQuery, id)

	var realEstate domain.RealEstate
	err := row.Scan(
		&realEstate.ID, &realEstate.UserID,
		&realEstate.Name, &realEstate.Currency,
		&realEstate.PurchasePriceCents, &realEstate.MonthlyIncomeCents,
		&realEstate.Purchased,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("get real estate err: %w", err)
	}

	return &realEstate, nil
}
