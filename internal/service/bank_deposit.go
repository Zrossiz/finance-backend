package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
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
	return b.pgBankDeposit.GetAllByUserID(ctx, userID)
}

func (b *bankDeposit) GetByID(ctx context.Context, id uuid.UUID) (*domain.BankDeposit, error) {
	return b.pgBankDeposit.GetByID(ctx, id)
}
