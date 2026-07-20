package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
)

type bankDeposit struct {
	conn *sql.DB
}

func newBankDeposit(conn *sql.DB) *bankDeposit {
	return &bankDeposit{conn: conn}
}

func (b *bankDeposit) Create(ctx context.Context, payload domain.BankDeposit) error {
	_, err := b.conn.ExecContext(
		ctx, createBankDepositQuery, payload.ID,
		payload.UserID, payload.Name, payload.Currency,
		payload.AmountCents, payload.InterestRate,
		payload.OpenedAt, payload.PeriodMonths, payload.TotalIncomeCents,
	)
	if err != nil {
		return fmt.Errorf("create bank deposit db err: %w", err)
	}

	return nil
}

func (b *bankDeposit) GetByID(ctx context.Context, id uuid.UUID) (*domain.BankDeposit, error) {
	row := b.conn.QueryRowContext(ctx, getBankDepositByIDQuery, id)

	var deposit domain.BankDeposit
	err := row.Scan(
		&deposit.ID, &deposit.UserID,
		&deposit.Name, &deposit.Currency,
		&deposit.AmountCents, &deposit.InterestRate,
		&deposit.OpenedAt, &deposit.PeriodMonths, &deposit.TotalIncomeCents,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("get bank deposit by id db err: %w", err)
	}

	return &deposit, nil
}

func (b *bankDeposit) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := b.conn.ExecContext(ctx, deleteBankDepositQuery, id)
	if err != nil {
		return fmt.Errorf("delete bank deposit err: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (b *bankDeposit) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankDeposit, error) {
	rows, err := b.conn.QueryContext(ctx, getUserBankDepositsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("get user bank deposits db err: %w", err)
	}

	var deposits []domain.BankDeposit
	for rows.Next() {
		var deposit domain.BankDeposit
		err := rows.Scan(
			&deposit.ID, &deposit.UserID,
			&deposit.Name, &deposit.Currency,
			&deposit.AmountCents, &deposit.InterestRate,
			&deposit.OpenedAt, &deposit.PeriodMonths, &deposit.TotalIncomeCents,
		)
		if err != nil {
			return nil, fmt.Errorf("scan bank deposit err: %w", err)
		}

		deposits = append(deposits, deposit)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate bank deposits err: %w", err)
	}

	return deposits, nil
}
