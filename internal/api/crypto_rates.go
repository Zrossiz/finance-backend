package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Zrossiz/finance-backend/internal/domain"
	"github.com/Zrossiz/finance-backend/internal/helpers"
)

type cryptoRates struct {
	apiURL     string
	httpClient *http.Client
}

func newCryptoRates() *cryptoRates {
	return &cryptoRates{
		apiURL: "https://api.coingecko.com/api/v3",
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *cryptoRates) GetByIds(ctx context.Context, ids []string) (domain.CryptoRates, error) {
	query := fmt.Sprintf(
		"%s/simple/price?ids=%s&vs_currencies=usd",
		c.apiURL,
		strings.Join(ids, ","),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, query, nil)
	if err != nil {
		return nil, fmt.Errorf("create crypto rates request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make crypto rates request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid crypto rates status code: %d", resp.StatusCode)
	}

	rates, err := helpers.ParseJSONBody[domain.CryptoRates](resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse crypto rates: %w", err)
	}

	return rates, nil
}
