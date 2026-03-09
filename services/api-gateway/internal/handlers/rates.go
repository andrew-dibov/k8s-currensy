package handlers

import (
	"api-gateway/internal/clients"
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

func (handler *RatesHandler) GetRates(res http.ResponseController, req *http.Request) {
	base := req.URL.Query().Get("base")

	if base == "" {
		base = "USD"
	}

	handler.logger.WithFields(logrus.Fields{
		"req_path":  req.URL.Path,
		"base_curr": base,
	}).Debug("getting rates")
}
