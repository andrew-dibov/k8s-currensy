package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ExchangeClient struct {
	baseURL string
	client  *http.Client
}

type ExchangeResponse struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func NewExchangeClient(url string, token string) *ExchangeClient {
	return &ExchangeClient{
		baseURL: fmt.Sprintf("%s%s/latest/", url, token),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

/* --- --- --- */

func (exchange *ExchangeClient) GetRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("%s%s", exchange.baseURL, baseCurrency)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request %v", err)
	}

	res, err := exchange.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates : %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with : %s", res.Status)
	}

	var data ExchangeResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response : %v", err)
	}

	if data.ConversionRates == nil {
		return nil, fmt.Errorf("API response miss rates")
	}

	return data.ConversionRates, nil
}
