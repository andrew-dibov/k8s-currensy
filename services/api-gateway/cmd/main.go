package main

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/configs"
	"api-gateway/internal/middlewares"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := configs.Load()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	ratesClient, err := clients.NewRatesClient(cfg.RatesService)
	if err != nil {
		logger.WithError(err).Fatal("failed to create rates client")
	}
	// defer ratesClient. // добавить метод close() в клиент

	// ratesHandler := handlers.NewRatesHandler(ratesClient, logger)

	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.LoggerMiddleware(next, logger)
	})

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.AuthenticatorMiddleware(next, logger, cfg.APIKeys)
	})

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("API-Gateway v1.0.0"))
	})

	router.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("OK"))
	})

	router.HandleFunc("/api/v1/rates", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	router.HandleFunc("/api/v1/convert", func(w http.ResponseWriter, r *http.Request) {}).Methods("POST")
	router.HandleFunc("/api/v1/convert", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

	http.ListenAndServe(":8080", router)

	logger.Info("starting API gateway on : ", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		logger.WithError(err).Fatal("server failed")
	}
}
