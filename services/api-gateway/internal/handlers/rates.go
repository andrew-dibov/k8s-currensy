package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"api-gateway/internal/clients"
)

type RatesHandler struct {
	ratesClient *clients.RatesClient
	logger      *logrus.Logger
}

func NewRatesHandler(ratesClient *clients.RatesClient, l *logrus.Logger) *RatesHandler {
	return &RatesHandler{
		ratesClient: ratesClient,
		logger:      l,
	}
}

func (h *RatesHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	h.logger.WithFields(logrus.Fields{
		"base_currency": baseCurrency,
		"path":          r.URL.Path,
	}).Debug("Processing rates")

	res, err := h.ratesClient.GetRates(r.Context(), baseCurrency)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get rates")
		http.Error(w, "Failed to fetch currency rates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.logger.WithError(err).Error("Failed to encode response to JSON")
	}
}
