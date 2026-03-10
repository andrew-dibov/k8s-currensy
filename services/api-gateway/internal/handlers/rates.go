package handlers

import (
	"api-gateway/internal/clients"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RatesHandler struct {
	client *clients.RatesClient
	logger *logrus.Logger
}

func NewRatesHandler(client *clients.RatesClient, logger *logrus.Logger) *RatesHandler {
	return &RatesHandler{
		client: client,
		logger: logger,
	}
}

func (h *RatesHandler) respond(w http.ResponseWriter, status int, res any) {
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

func (h *RatesHandler) GetRate(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from_currency")
	to := r.URL.Query().Get("to_currency")

	if from == "" || to == "" {
		http.Error(w, "missing params", http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	data, err := h.client.GetRate(ctx, from, to)
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

func (h *RatesHandler) GetAllRates(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base_currency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	ctx := r.Context()

	data, err := h.client.GetAllRates(ctx, baseCurrency)
	if err != nil {
		h.logger.WithError(err).Error("failed to get rates")
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	h.respond(w, http.StatusOK, GetAllRatesResponse{
		BaseCurrency: data.BaseCurrency,
		Rates:        data.Rates,
	})
}

/* --- --- --- */

type ConvertRequest struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
}

type ConvertResponse struct {
	Result       float64 `json:"result"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Rate         float64 `json:"rate"`
	Amount       float64 `json:"amount"`
}

func (h *RatesHandler) Convert(w http.ResponseWriter, r *http.Request) {
	var body ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.FromCurrency == "" || body.ToCurrency == "" || body.Amount <= 0 {
		http.Error(w, "invalid fields", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	data, err := h.client.Convert(ctx, body.FromCurrency, body.ToCurrency, body.Amount)
	if err != nil {
		h.logger.WithError(err).Error("failed to convert")
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	h.respond(w, http.StatusOK, ConvertResponse{
		Result:       data.Result,
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Rate:         data.Rate,
		Amount:       data.Amount,
	})
}
