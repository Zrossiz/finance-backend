package service

import (
	"context"

	"github.com/Zrossiz/finance-backend/internal/config"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type IStockRepo interface {
	Create(ctx context.Context, payload domain.Stock) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Stock, error)
	GetOne(ctx context.Context, id uuid.UUID) (*domain.Stock, error)
	Update(ctx context.Context, payload domain.Stock) error
}

type IBondRepo interface {
	Create(ctx context.Context, payload domain.Bond) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Bond, error)
	GetOne(ctx context.Context, id uuid.UUID) (*domain.Bond, error)
}

type ICryptoPositionRepo interface {
	Create(ctx context.Context, payload domain.CryptoPosition) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CryptoPosition, error)
	GetOneByID(ctx context.Context, id uuid.UUID) (*domain.CryptoPosition, error)
	Update(ctx context.Context, id uuid.UUID, amount decimal.Decimal, avgPriceUsd int64) error
}

type IRealEstateRepo interface {
	Create(ctx context.Context, payload domain.RealEstate) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.RealEstate, error)
	GetOne(ctx context.Context, id uuid.UUID) (*domain.RealEstate, error)
	Update(ctx context.Context, payload domain.RealEstate) error
}

type IUserRepo interface {
	Create(ctx context.Context, payload domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type IBankDepositRepo interface {
	Create(ctx context.Context, payload domain.BankDeposit) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankDeposit, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BankDeposit, error)
}

type Service struct {
	User           *user
	Stock          *stock
	BankDeposit    *bankDeposit
	Bond           *bond
	RealEstate     *realEstate
	CryptoPosition *cryptoPosition
}

type Postgres struct {
	User           IUserRepo
	Stock          IStockRepo
	BankDeposit    IBankDepositRepo
	Bond           IBondRepo
	RealEstate     IRealEstateRepo
	CryptoPosition ICryptoPositionRepo
}

func NewService(pgRepo Postgres, cfg *config.Config) *Service {
	return &Service{
		User:           newUser(pgRepo.User, cfg.Server.JWTAccessSecret, cfg.Server.JWTRefreshSecret),
		Stock:          newStock(pgRepo.Stock),
		BankDeposit:    newBankDeposit(pgRepo.BankDeposit),
		Bond:           newBond(pgRepo.Bond),
		RealEstate:     newRealEstate(pgRepo.RealEstate),
		CryptoPosition: newCryptoPosition(pgRepo.CryptoPosition),
	}
}
