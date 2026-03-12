package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ExchangeClient struct {
	url    string
	client *http.Client
}

type ExchangeRes struct {
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

func NewExchangeClient(url string, token string) *ExchangeClient {
	return &ExchangeClient{
		url: fmt.Sprintf("%s%s/latest/", url, token),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (ec *ExchangeClient) GetRates(baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("%s%s", ec.url, baseCurrency)

	if baseCurrency == "" {
		return nil, fmt.Errorf("empty base currency")
	}

	res, err := ec.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates : %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with : %s", res.Status)
	}

	var data ExchangeRes
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response : %v", err)
	}

	return data.ConversionRates, nil
}
