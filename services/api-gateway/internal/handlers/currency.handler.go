package handlers

import (
	"api-gateway/internal/clients"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type CurrencyHandler struct {
	client             *clients.CurrencyClient
	logger             *logrus.Logger
	getRateLimiter     *rate.Limiter
	getAllRatesLimiter *rate.Limiter
}

func NewCurrencyHandler(client *clients.CurrencyClient, logger *logrus.Logger) *CurrencyHandler {
	return &CurrencyHandler{
		client:             client,
		logger:             logger,
		getRateLimiter:     rate.NewLimiter(rate.Limit(10), 30),
		getAllRatesLimiter: rate.NewLimiter(rate.Limit(10), 30),
	}
}

func (handler *CurrencyHandler) respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		handler.logger.WithError(err).Error("failed to encode JSON")
	}
}

func (handler *CurrencyHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		handler.logger.WithError(err).Error("failed to encode JSON")
	}
}

/* --- --- --- */

type GetRateResponse struct {
	Rate         float64 `json:"rate"`
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
}

func (handler *CurrencyHandler) GetRate(w http.ResponseWriter, r *http.Request) {
	if !handler.getRateLimiter.Allow() {
		handler.logger.WithField("remote_addr", r.RemoteAddr).Warn("rate limit exceeded")
		handler.respondError(w, http.StatusTooManyRequests, "too many requests")
		return
	}

	fromCurrency := r.URL.Query().Get("fromCurrency")
	toCurrency := r.URL.Query().Get("toCurrency")

	if fromCurrency == "" || toCurrency == "" {
		handler.respondError(w, http.StatusBadRequest, "fromCurrency or toCurrency is empty")
		return
	}

	if len(fromCurrency) != 3 || len(toCurrency) != 3 {
		handler.respondError(w, http.StatusBadRequest, "fromCurrency or toCurrency length is not 3 chars")
		return
	}

	for _, ch := range fromCurrency {
		if ch < 'A' || ch > 'Z' {
			handler.respondError(w, http.StatusBadRequest, "fromCurrency has invalid chars")
			return
		}
	}

	for _, ch := range toCurrency {
		if ch < 'A' || ch > 'Z' {
			handler.respondError(w, http.StatusBadRequest, "toCurrency has invalid chars")
			return
		}
	}

	ctx := r.Context()

	data, err := handler.client.GetRate(
		ctx,
		fromCurrency,
		toCurrency,
	)

	if err != nil {
		handler.logger.WithError(err).Error("get rate failed")
		handler.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	handler.respond(w, http.StatusOK, GetRateResponse{
		Rate:         data.Rate,
		FromCurrency: data.FromCurrency,
		ToCurrency:   data.ToCurrency,
	})
}

type GetAllRatesResponse struct {
	BaseCurrency string             `json:"baseCurrency"`
	Rates        map[string]float64 `json:"rates"`
}

func (handler *CurrencyHandler) GetAllRates(w http.ResponseWriter, r *http.Request) {
	if !handler.getAllRatesLimiter.Allow() {
		handler.logger.WithField("remote_addr", r.RemoteAddr).Warn("rate limit exceeded")
		handler.respondError(w, http.StatusTooManyRequests, "too many requests")
		return
	}

	baseCurrency := r.URL.Query().Get("baseCurrency")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	for _, ch := range baseCurrency {
		if ch < 'A' || ch > 'Z' {
			handler.respondError(w, http.StatusBadRequest, "baseCurrency has invalid chars")
			return
		}
	}

	if len(baseCurrency) != 3 {
		handler.respondError(w, http.StatusBadRequest, "baseCurrency length is not 3 chars")
		return
	}

	ctx := r.Context()

	data, err := handler.client.GetAllRates(
		ctx,
		baseCurrency,
	)

	if err != nil {
		handler.logger.WithError(err).Error("get all rates failed")
		handler.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	handler.respond(w, http.StatusOK, GetAllRatesResponse{
		BaseCurrency: data.BaseCurrency,
		Rates:        data.Rates,
	})
}
