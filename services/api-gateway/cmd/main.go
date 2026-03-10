package main

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/configs"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middlewares"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := configs.Load()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	/* --- --- --- */

	ratesClient, err := clients.NewRatesClient(cfg.RatesServiceURL)
	if err != nil {
		log.WithError(err).Fatal("failed to init rates client")
	}
	defer ratesClient.Close()

	rates := handlers.NewRatesHandler(ratesClient, log)

	/* --- --- --- */

	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.LoggerMiddleware(next, log)
	})

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.AuthenticatorMiddleware(next, log, cfg.APIKeys)
	})

	/* --- --- --- */

	router.HandleFunc("/", handlers.RootHandler)
	router.HandleFunc("/health", handlers.HealthHandler)

	router.HandleFunc("/api/v1/rate", rates.GetRate).Methods("GET")
	router.HandleFunc("/api/v1/all_rates", rates.GetAllRates).Methods("GET")
	router.HandleFunc("/api/v1/convert", rates.Convert).Methods("POST")

	/* --- --- --- */

	http.ListenAndServe(":8080", router)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.WithError(err).Fatal("server failed")
	}
}
