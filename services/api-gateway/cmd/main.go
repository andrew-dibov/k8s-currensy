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

	currencyClient, err := clients.NewCurrencyClient(cfg.CurrencyService)
	if err != nil {
		log.WithError(err).Fatal("failed to init currency client")
	}
	defer currencyClient.Close()

	// conversionClient, err := clients.NewConversionClient(cfg.ConversionService)
	// if err != nil {
	// 	log.WithError(err).Fatal("failed to init conversion client")
	// }
	// defer conversionClient.Close()

	/* --- --- --- */

	currency := handlers.NewCurrencyHandler(currencyClient, log)
	// conversion := handlers.NewConversionHandler(conversionClient, log)

	/* --- --- --- */

	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.LoggerMiddleware(next, log)
	})

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.AuthenticatorMiddleware(next, log, cfg.APIKeys)
	})

	/* --- --- --- */

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/health", healthHandler)

	router.HandleFunc("/api/v1/rate", currency.GetRate).Methods("GET")
	router.HandleFunc("/api/v1/allRates", currency.GetAllRates).Methods("GET")
	// router.HandleFunc("/api/v1/convert", conversion.Convert).Methods("POST")

	/* --- --- --- */

	http.ListenAndServe(":8080", router)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.WithError(err).Fatal("server failed")
	}
}

/* --- --- --- */

func rootHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("API-Gateway v1.0.0"))
}

func healthHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK"))
}
