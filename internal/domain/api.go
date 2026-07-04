package domain

type CryptoRates map[string]CryptoRate

type CryptoRate struct {
	USD string `json:"usd"`
}
