package handlers

import (
	"api-gateway/internal/clients"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type ConversionHandler struct {
	client         *clients.ConversionClient
	logger         *logrus.Logger
	convertLimiter *rate.Limiter
}

func NewConversionHandler(client *clients.ConversionClient, logger *logrus.Logger) *ConversionHandler {
	return &ConversionHandler{
		client:         client,
		logger:         logger,
		convertLimiter: rate.NewLimiter(rate.Limit(10), 30),
	}
}

func (handler *ConversionHandler) respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		handler.logger.WithError(err).Error("failed to encode JSON")
	}
}

func (handler *ConversionHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		handler.logger.WithError(err).Error("failed to encode JSON")
	}
}

/* --- --- --- */

type ConvertRequest struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Amount       float64 `json:"amount"`
}

type ConvertResponse struct {
	Result       float64 `json:"result"`
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Rate         float64 `json:"rate"`
	Amount       float64 `json:"amount"`
}

func (handler *ConversionHandler) Convert(w http.ResponseWriter, r *http.Request) {
	if !handler.convertLimiter.Allow() {
		handler.logger.WithField("remote_addr", r.RemoteAddr).Warn("rate limit exceeded")
		handler.respondError(w, http.StatusTooManyRequests, "too many requests")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	var body ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handler.respondError(w, http.StatusBadRequest, "body is invalid")
		return
	}

	if body.FromCurrency == "" || body.ToCurrency == "" {
		handler.respondError(w, http.StatusBadRequest, "fromCurrency or toCurrency is empty")
		return
	}

	if len(body.FromCurrency) != 3 || len(body.ToCurrency) != 3 {
		handler.respondError(w, http.StatusBadRequest, "fromCurrency or toCurrency length is not 3 chars")
		return
	}

	for _, ch := range body.FromCurrency {
		if ch < 'A' || ch > 'Z' {
			handler.respondError(w, http.StatusBadRequest, "fromCurrency has invalid chars")
			return
		}
	}

	for _, ch := range body.ToCurrency {
		if ch < 'A' || ch > 'Z' {
			handler.respondError(w, http.StatusBadRequest, "toCurrency has invalid chars")
			return
		}
	}

	if body.Amount <= 0 {
		handler.respondError(w, http.StatusBadRequest, "amount must be positive")
		return
	}

	if body.Amount > 1e12 {
		handler.respondError(w, http.StatusBadRequest, "amount is too large")
		return
	}

	ctx := r.Context()

	data, err := handler.client.Convert(
		ctx,
		body.FromCurrency,
		body.ToCurrency,
		body.Amount,
	)

	if err != nil {
		handler.logger.WithError(err).Error("conversion failed")
		handler.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	handler.respond(w, http.StatusOK, ConvertResponse{
		Result:       data.Result,
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
		Rate:         data.Rate,
		Amount:       data.Amount,
	})
}
