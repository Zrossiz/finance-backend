package handler

import (
	"time"

	"github.com/Zrossiz/finance-backend/internal/domain"
)

type createUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createCryptoPositionDTO struct {
	Ticker           string `json:"ticker"`
	Amount           string `json:"amount"`
	CoinID           string `json:"coin_id"`
	AvgPriceUSDCents *int64 `json:"avg_price_usd_cents"`
}

type updateCryptoPositionDTO struct {
	Amount           string `json:"amount"`
	AvgPriceUSDCents *int64 `json:"avg_price_usd_cents"`
}

type getUserCryptoPositionsResDTO struct {
	Total       string                  `json:"total"`
	TotalProfit string                  `json:"total_profit"`
	Positions   []domain.CryptoPosition `json:"positions"`
}

type getCryptoCoinsResDTO struct {
	Coins []domain.CryptoCoin `json:"coins"`
}

type createBankDepositDTO struct {
	Name         string    `json:"name"`
	Currency     string    `json:"currency"`
	AmountCents  int64     `json:"amount_cents"`
	InterestRate string    `json:"interest_rate"`
	OpenedAt     time.Time `json:"opened_at"`
	PeriodMonths int       `json:"period_months"`
}

type createRealEstateDTO struct {
	Name               string     `json:"name"`
	Currency           string     `json:"currency"`
	PurchasePriceCents *int64     `json:"purchase_price_cents"`
	MonthlyIncomeCents *int64     `json:"monthly_income_cents"`
	Purchased          *time.Time `json:"purchased"`
}

type updateRealEstateDTO struct {
	Name               string     `json:"name"`
	Currency           string     `json:"currency"`
	PurchasePriceCents *int64     `json:"purchase_price_cents"`
	MonthlyIncomeCents *int64     `json:"monthly_income_cents"`
	Purchased          *time.Time `json:"purchased"`
}
