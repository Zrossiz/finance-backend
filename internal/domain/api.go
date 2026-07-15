package domain

type CryptoRates map[string]CryptoRate

type CryptoRate struct {
	USD float32 `json:"usd"`
}
