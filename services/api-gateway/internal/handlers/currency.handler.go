package handlers

import (
	"api-gateway/internal/clients"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type CurrencyHandler struct {
	client *clients.CurrencyClient
	logger *logrus.Logger
}

func NewCurrencyHandler(c *clients.CurrencyClient, l *logrus.Logger) *CurrencyHandler {
	return &CurrencyHandler{
		client: c,
		logger: l,
	}
}

func (h *CurrencyHandler) respond(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.logger.WithError(err).Error("failed to encode JSON")
	}
}

/* --- --- --- */

type GetRateResponse struct {
	Rate         float64 `json:"rate"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
}

func (h *CurrencyHandler) GetRate(w http.ResponseWriter, r *http.Request) {
	fromCurrency := r.URL.Query().Get("from_currency")
	toCurrency := r.URL.Query().Get("to_currency")

	if fromCurrency == "" || toCurrency == "" {
		http.Error(w, "missing params", http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	data, err := h.client.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		h.logger.WithError(err).Error("failed to get rate")
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	h.respond(w, http.StatusOK, GetRateResponse{
		Rate:         data.Rate,
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
	})

}

/* --- --- --- */

type GetAllRatesResponse struct {
	BaseCurrency string             `json:"base_currency"`
	Rates        map[string]float64 `json:"rates"`
}

func (h *CurrencyHandler) GetAllRates(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base_currency")

	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	ctx := r.Context()

	data, err := h.client.GetAllRates(ctx, baseCurrency)
	if err != nil {
		h.logger.WithError(err).Error("failed to get all rates")
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	h.respond(w, http.StatusOK, GetAllRatesResponse{
		BaseCurrency: data.BaseCurrency,
		Rates:        data.Rates,
	})
}
