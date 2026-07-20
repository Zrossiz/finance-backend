package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type bankDeposit struct {
	pgBankDeposit IBankDepositRepo
}

func newBankDeposit(pgBankDeposit IBankDepositRepo) *bankDeposit {
	return &bankDeposit{pgBankDeposit: pgBankDeposit}
}

func (b *bankDeposit) Create(ctx context.Context, payload domain.BankDeposit) error {
	return b.pgBankDeposit.Create(ctx, payload)
}

func (b *bankDeposit) Delete(ctx context.Context, id uuid.UUID) error {
	return b.pgBankDeposit.Delete(ctx, id)
}

func (b *bankDeposit) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankDeposit, error) {
	deposits, err := b.pgBankDeposit.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range deposits {
		deposits[i].IncomeCentsPerMonth = deposits[i].TotalIncomeCents / int64(deposits[i].PeriodMonths)
	}

	return deposits, nil
}

func (b *bankDeposit) GetByID(ctx context.Context, id uuid.UUID) (*domain.BankDeposit, error) {
	return b.pgBankDeposit.GetByID(ctx, id)
}

func (b *bankDeposit) CalculateTotalIncomeCents(
	amountCents int64,
	periodMonths int,
	interestRate decimal.Decimal,
) int64 {
	return decimal.
		NewFromInt(amountCents).
		Mul(interestRate).
		Div(decimal.NewFromInt(100)).
		Mul(decimal.NewFromInt(int64(periodMonths))).
		Div(decimal.NewFromInt(12)).
		IntPart()
}
