package handler

import (
	"context"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/config"
	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type IStockService interface {
	Create(ctx context.Context, payload domain.Stock) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Stock, error)
	GetOne(ctx context.Context, id uuid.UUID) (*domain.Stock, error)
	Update(ctx context.Context, payload domain.Stock) error
}

type IBondService interface {
	Create(ctx context.Context, payload domain.Bond) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Bond, error)
	GetOne(ctx context.Context, id uuid.UUID) (*domain.Bond, error)
}

type ICryptoPositionService interface {
	Create(ctx context.Context, payload domain.CryptoPosition) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.CryptoPosition, error)
	GetOneByID(ctx context.Context, id uuid.UUID) (*domain.CryptoPosition, error)
	Update(ctx context.Context, id uuid.UUID, amount decimal.Decimal, avgPriceUsdCents *int64) error
	CountTotalByPositions(positions []domain.CryptoPosition) decimal.Decimal
	CountTotalProfitByPositions(positions []domain.CryptoPosition) decimal.Decimal
}

type IRealEstateService interface {
	Create(ctx context.Context, payload domain.RealEstate) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.RealEstate, error)
	GetOne(ctx context.Context, id uuid.UUID) (*domain.RealEstate, error)
	Update(ctx context.Context, payload domain.RealEstate) error
}

type IUserService interface {
	Login(ctx context.Context, username, password string) (string, string, error)
	Registration(ctx context.Context, username, password string) (string, string, error)
	RefreshTokens(ctx context.Context, refresh string) (string, string, error)
}

type IBankDepositService interface {
	Create(ctx context.Context, payload domain.BankDeposit) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BankDeposit, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BankDeposit, error)
	CalculateTotalIncomeCents(amountCents int64, period int, interestRate decimal.Decimal) int64
}

type ICryptoCoinService interface {
	GetAll(ctx context.Context) ([]domain.CryptoCoin, error)
}

type Service struct {
	User           IUserService
	Stock          IStockService
	BankDeposit    IBankDepositService
	Bond           IBondService
	RealEstate     IRealEstateService
	CryptoPosition ICryptoPositionService
	CryptoCoin     ICryptoCoinService
}

type Handler struct {
	user           *user
	cryptoPosition *cryptoPosition
	cryptoCoin     *cryptoCoin
	bankDeposit    *bankDeposit
}

func New(srv Service, cfg *config.Config) (*Handler, error) {
	userHandler, err := newUser(srv.User, cfg)
	if err != nil {
		return nil, fmt.Errorf("init user handler err: %w", err)
	}

	return &Handler{
		user:           userHandler,
		cryptoPosition: newCryptoPosition(srv.CryptoPosition),
		cryptoCoin:     newCryptoCoin(srv.CryptoCoin),
		bankDeposit:    newBankDeposit(srv.BankDeposit),
	}, nil
}

func (h *Handler) RegisterRoutes(router chi.Router, accessSecret string) {
	router.Route("/api/v1", func(r chi.Router) {
		h.user.registerRoutes(r)
		h.cryptoCoin.registerRoutes(r)
		h.bankDeposit.registerRoutes(r, accessSecret)
		h.cryptoPosition.registerRoutes(r, accessSecret)
	})
}
