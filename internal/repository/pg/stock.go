package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
)

type stock struct {
	conn *sql.DB
}

func newStock(conn *sql.DB) *stock {
	return &stock{conn: conn}
}

func (s *stock) Create(ctx context.Context, payload domain.Stock) error {
	_, err := s.conn.ExecContext(
		ctx, createStockQuery,
		payload.ID, payload.UserID,
		payload.Ticker, payload.Currency,
		payload.Amount, payload.AvgPriceCents,
	)
	if err != nil {
		return fmt.Errorf("create stock error: %w", err)
	}

	return nil
}

func (s *stock) GetOne(ctx context.Context, id uuid.UUID) (*domain.Stock, error) {
	row := s.conn.QueryRowContext(ctx, getStockQuery, id)

	var stock domain.Stock
	err := row.Scan(
		&stock.ID, &stock.UserID,
		&stock.Ticker, &stock.Currency,
		&stock.Amount, &stock.AvgPriceCents,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("get stock db err: %w", err)
	}

	return &stock, nil
}

func (s *stock) Update(ctx context.Context, payload domain.Stock) error {
	res, err := s.conn.ExecContext(
		ctx, updateStockQuery,
		payload.ID, payload.UserID,
		payload.Ticker, payload.Currency,
		payload.Amount, payload.AvgPriceCents,
	)
	if err != nil {
		return fmt.Errorf("update stock error: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (s *stock) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := s.conn.ExecContext(ctx, deleteStockQuery, id)
	if err != nil {
		return fmt.Errorf("delete stock db err: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (s *stock) GetAllByUserID(
	ctx context.Context, userID uuid.UUID,
) ([]domain.Stock, error) {
	rows, err := s.conn.QueryContext(ctx, getUserStocksQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("get user stocks db err: %w", err)
	}

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(
			&stock.ID, &stock.UserID,
			&stock.Ticker, &stock.Currency,
			&stock.Amount, &stock.AvgPriceCents,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user sotck db err: %w", err)
		}

		stocks = append(stocks, stock)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate user stocks db err: %w", err)
	}

	return stocks, nil
}
