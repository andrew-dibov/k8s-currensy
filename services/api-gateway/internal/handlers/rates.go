package handlers

import (
	"api-gateway/internal/clients"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RatesHandler struct {
	ratesClient *clients.RatesClient
	logger      *logrus.Logger
}

func NewRatesHandler(ratesClient *clients.RatesClient, logger *logrus.Logger) *RatesHandler {
	return &RatesHandler{
		ratesClient: ratesClient,
		logger:      logger,
	}
}

func (h *RatesHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base")
	/* YOOOOOOOOOOOOOOOOOOOO */
}
