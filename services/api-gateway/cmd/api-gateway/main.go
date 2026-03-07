package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"api-gateway/internal/middlewares"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	router := mux.NewRouter()

	/*
		Middleware : логирование запросов
			1. Use(mwf) : метод маршрутизатора : добавить мидлвар в цепь
			2. func(next) : анонимная функция : принять следующий в цепи обработчик
			3. LoggerMiddleware(next, logger) : передать обработчик + логгер в мидлвар
	*/

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.LoggerMiddleware(next, logger)
	})

	/*
		Middleware : аутентификация запросов
			1. Use(mwf) : метод маршрутизатора : добавить мидлвар в цепь
			2. func(next) : анонимная функция : принять следующий в цепи обработчик
			3. AuthenticatorMiddleware(next, logger) : передать обработчик + логгер в мидлвар
	*/

	router.Use(func(next http.Handler) http.Handler {
		return middlewares.AuthenticatorMiddleware(next, logger)
	})

	/* HandleFunc(path, func) : обработка запросов : путь -> функция */

	router.HandleFunc("/", appHandler)
	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/api/v1/rates", ratesHandler)
	router.HandleFunc("/api/v1/convert", convertHandler)
	router.HandleFunc("/api/v1/history", historyHandler)

	/* ListenAndServe(port, router) : запуск сервера : порт -> маршрутизатор */

	http.ListenAndServe(":8080", router)
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API Gateway v1.0.0")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "API Gateway is up")
}

func ratesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Rates handler")
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Convert handler")
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "History handler")
}
