package fetchers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ExchangeRateFetcher struct {
	url    string
	client *http.Client
}

type ExchangeRateResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

func NewExchangeRateFetcher(url string, token string) *ExchangeRateFetcher {
	return &ExchangeRateFetcher{
		url: fmt.Sprintf("%s%s/latest/", url, token),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (f *ExchangeRateFetcher) FetchRates(baseCurrency string) (map[string]float64, error) {
	reqURL := fmt.Sprintf("%s%s", f.url, baseCurrency)

	if baseCurrency == "" {
		return nil, fmt.Errorf("base currency is empty")
	}

	res, err := f.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates : %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned : %s", res.Status)
	}

	var data ExchangeRateResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response : %v", err)
	}

	return data.ConversionRates, nil
}
