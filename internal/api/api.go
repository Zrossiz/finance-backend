package api

type API struct {
	CryptoRates *cryptoRates
}

func NewAPI() *API {
	return &API{
		CryptoRates: newCryptoRates(),
	}
}
