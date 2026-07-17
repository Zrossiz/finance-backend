package service

import (
	"context"
	"fmt"

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
	Update(ctx context.Context, id uuid.UUID, amount decimal.Decimal, avgPriceUsd *int64) error
	GetUniqueCoinIDs(ctx context.Context) ([]string, error)
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

type ICryptoCoinRepo interface {
	GetAll(ctx context.Context) ([]domain.CryptoCoin, error)
}

type ICryptoRatesAPI interface {
	GetByIDs(ctx context.Context, ids []string) (domain.CryptoRates, error)
}

type ICryptoRatesCache interface {
	GetByIDs(ctx context.Context, ids []string) (domain.CryptoRates, error)
	Save(ctx context.Context, rates domain.CryptoRates) error
}

type Postgres struct {
	User           IUserRepo
	Stock          IStockRepo
	BankDeposit    IBankDepositRepo
	Bond           IBondRepo
	RealEstate     IRealEstateRepo
	CryptoPosition ICryptoPositionRepo
	CryptoCoin     ICryptoCoinRepo
}

type API struct {
	CryptoRates ICryptoRatesAPI
}

type Cache struct {
	CryptoRates ICryptoRatesCache
}

type Service struct {
	User           *user
	Stock          *stock
	BankDeposit    *bankDeposit
	Bond           *bond
	RealEstate     *realEstate
	CryptoPosition *cryptoPosition
	CryptoRates    *cryptoRate
	CryptoCoin     *cryptoCoin
}

func New(pgRepo Postgres, apiSrv API, cache Cache, cfg *config.Config) (*Service, error) {
	userSrv, err := newUser(pgRepo.User, cfg)
	if err != nil {
		return nil, fmt.Errorf("init user service err: %w", err)
	}

	return &Service{
		User:           userSrv,
		Stock:          newStock(pgRepo.Stock),
		BankDeposit:    newBankDeposit(pgRepo.BankDeposit),
		Bond:           newBond(pgRepo.Bond),
		RealEstate:     newRealEstate(pgRepo.RealEstate),
		CryptoPosition: newCryptoPosition(pgRepo.CryptoPosition, apiSrv.CryptoRates, cache.CryptoRates),
		CryptoRates:    newCryptoRate(pgRepo.CryptoCoin, apiSrv.CryptoRates, cache.CryptoRates),
		CryptoCoin:     newCryptoCoin(pgRepo.CryptoCoin),
	}, nil
}
