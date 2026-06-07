package main

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/configs"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middlewares"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	appConfig := configs.LoadAppConfig()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"app_port": appConfig.AppPort,
		"app_keys": appConfig.AppKeys,

		"history_url":    appConfig.HistoryURL,
		"currency_url":   appConfig.CurrencyURL,
		"conversion_url": appConfig.ConversionURL,
	}).Info("api-gateway : app config loaded")

	/* --- --- --- */

	currencyClient, err := clients.NewCurrencyClient(appConfig.CurrencyURL, 5*time.Second)
	if err != nil {
		logger.WithError(err).Fatal("currency client failed")
	}
	defer currencyClient.Close()

	currency := handlers.NewCurrencyHandler(currencyClient, logger)

	conversionClient, err := clients.NewConversionClient(appConfig.ConversionURL, 5*time.Second)
	if err != nil {
		logger.WithError(err).Fatal("currency client failed")
	}
	defer conversionClient.Close()

	conversion := handlers.NewConversionHandler(conversionClient, logger)

	/* --- --- --- */

	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.LoggerMiddleware(next, logger)
	})

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.AuthMiddleware(next, logger, appConfig.AppKeys)
	})

	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	router.HandleFunc("/api/v1/rate", currency.GetRate).Methods("GET")
	router.HandleFunc("/api/v1/rates", currency.GetAllRates).Methods("GET")
	router.HandleFunc("/api/v1/convert", conversion.Convert).Methods("POST")

	/* --- --- --- */

	server := &http.Server{
		Addr:         appConfig.AppPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.WithField("port", appConfig.AppPort).Info("server starting")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatal("server failed")
		}
	}()

	/* --- --- --- */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("server forced to stop")
	}

	logger.Info("server stopped")
}
