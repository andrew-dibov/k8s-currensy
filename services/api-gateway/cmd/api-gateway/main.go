package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"api-gateway/internal/middlewares"
)

// ---

type ProxyClient struct {
	client *http.Client
}

func NewProxyClient(t time.Duration) *ProxyClient {
	return &ProxyClient{
		client: &http.Client{
			Timeout: t,
		},
	}
}

func (p *ProxyClient) ForwardRequest(method string, url string, body any, headers http.Header) (*http.Response, error) {
	var reader io.Reader

	if body != nil {
		json, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	return p.client.Do(req)
}

// ---

func main() {
	l := logrus.New()
	r := mux.NewRouter()

	l.SetFormatter(&logrus.JSONFormatter{})

	r.Use(func(h http.Handler) http.Handler {
		return middlewares.LoggerMiddleware(h, l)
	})

	r.Use(func(h http.Handler) http.Handler {
		return middlewares.AuthenticatorMiddleware(h, l)
	})

	r.HandleFunc("/", appHandler)
	r.HandleFunc("/healthz", healthzHandler)

	r.HandleFunc("/api/v1/rates", ratesHandler)
	r.HandleFunc("/api/v1/convert", convertHandler)
	r.HandleFunc("/api/v1/history", historyHandler)

	l.Info("Server up")
	http.ListenAndServe(":8080", r)
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API Gateway v1.0.0")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
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
