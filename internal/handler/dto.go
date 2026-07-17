package handler

import "github.com/Zrossiz/finance-backend/internal/domain"

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
