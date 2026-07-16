package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type cryptoCoin struct {
	cryptoCoinSrv ICryptoCoinService
}

func newCryptoCoin(cryptoCoinSrv ICryptoCoinService) *cryptoCoin {
	return &cryptoCoin{cryptoCoinSrv: cryptoCoinSrv}
}

func (c *cryptoCoin) registerRoutes(r chi.Router) {
	r.Get("/crypto-coins", c.getAll)
}

func (c *cryptoCoin) getAll(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	coins, err := c.cryptoCoinSrv.GetAll(ctx)
	if err != nil {
		logrus.Errorf("get crypto coins err: %v", err)
		Error(rw, ErrInternalServerError)
		return
	}

	res := getCryptoCoinsResDTO{
		Coins: coins,
	}

	JSON(rw, http.StatusOK, res)
}
