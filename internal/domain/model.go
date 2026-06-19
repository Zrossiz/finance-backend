package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CryptoPosition struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Ticker           string
	Amount           decimal.Decimal
	AvgPriceUSDCents *int64
}

type BankDeposit struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Name         string
	Currency     string
	AmountCents  int64
	InterestRate decimal.Decimal
	OpenedAt     time.Time
	ClosedAt     *time.Time
}

type RealEstate struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Name               string
	Currency           string
	PurchasePriceCents *int64
	MonthlyIncomeCents *int64
	Purchased          *time.Time
}

type Stock struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Ticker        string
	Currency      string
	Amount        int64
	AvgPriceCents int64
}

type Bond struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Ticker             string
	Currency           string
	Amount             int64
	AvgPriceCents      int64
	CouponCents        int64
	CouponPeriodMonths int
}
