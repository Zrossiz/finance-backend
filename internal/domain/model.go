package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CryptoPosition struct {
	ID               uuid.UUID       `json:"id"`
	UserID           uuid.UUID       `json:"user_id"`
	Ticker           string          `json:"ticker"`
	CoinID           string          `json:"coin_id"`
	Amount           decimal.Decimal `json:"amount"`
	TotalPriceUSD    decimal.Decimal `json:"total_price_usd"`
	AvgPriceUSDCents *int64          `json:"avg_price_usd_cents"`
	ProfitUSD        decimal.Decimal `json:"profit_usd"`
}

type BankDeposit struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	Name         string          `json:"name"`
	Currency     string          `json:"currency"`
	AmountCents  int64           `json:"amount_cents"`
	InterestRate decimal.Decimal `json:"interest_rate"`
	OpenedAt     time.Time       `json:"opened_at"`
	ClosedAt     *time.Time      `json:"closed_at"`
}

type RealEstate struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	Name               string     `json:"name"`
	Currency           string     `json:"currency"`
	PurchasePriceCents *int64     `json:"purchase_price_cents"`
	MonthlyIncomeCents *int64     `json:"monthly_income_cents"`
	Purchased          *time.Time `json:"purchased"`
}

type Stock struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Ticker        string    `json:"ticker"`
	Currency      string    `json:"currency"`
	Amount        int64     `json:"amount"`
	AvgPriceCents int64     `json:"avg_price_cents"`
}

type Bond struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	Ticker             string    `json:"ticker"`
	Currency           string    `json:"currency"`
	Amount             int64     `json:"amount"`
	AvgPriceCents      int64     `json:"avg_price_cents"`
	CouponCents        int64     `json:"coupon_cents"`
	CouponPeriodMonths int       `json:"coupon_period_months"`
}

type CryptoCoin struct {
	Symbol string `json:"symbol"`
	CoinID string `json:"coin_id"`
}
