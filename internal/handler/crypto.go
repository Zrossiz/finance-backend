package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type cryptoPosition struct {
	cryptoPositionSrv ICryptoPositionService
}

func newCryptoPosition(cryptoPositionSrv ICryptoPositionService) *cryptoPosition {
	return &cryptoPosition{cryptoPositionSrv: cryptoPositionSrv}
}

func (c *cryptoPosition) registerRoutes(r chi.Router) {
	r.Get("/cryptos", c.getAllByUserID)
	r.Post("/cryptos", c.create)
	r.Delete("/cryptos/{id}", c.delete)
	r.Patch("/cryptos/{id}", c.update)
}

func (c *cryptoPosition) getAllByUserID(rw http.ResponseWriter, r *http.Request) {

}

func (c *cryptoPosition) create(rw http.ResponseWriter, r *http.Request) {

}

func (c *cryptoPosition) update(rw http.ResponseWriter, r *http.Request) {

}

func (c *cryptoPosition) delete(rw http.ResponseWriter, r *http.Request) {

}
