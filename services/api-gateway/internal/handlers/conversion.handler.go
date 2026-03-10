package handlers

import (
	"api-gateway/internal/clients"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ConversionHandler struct {
	client *clients.ConversionClient
	logger *logrus.Logger
}

func NewConversionHandler(c *clients.ConversionClient, l *logrus.Logger) *ConversionHandler {
	return &ConversionHandler{
		client: c,
		logger: l,
	}
}

func (h *ConversionHandler) respond(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.logger.WithError(err).Error("failed to encode JSON")
	}
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

func (h *ConversionHandler) Convert(w http.ResponseWriter, r *http.Request) {
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
