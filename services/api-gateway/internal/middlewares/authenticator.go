package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func AuthenticatorMiddleware(next http.Handler, logger *logrus.Logger, validKeys map[string]bool) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/health" {
			next.ServeHTTP(res, req)
			return
		}

		key := req.Header.Get("X-API-Key")

		if key == "" {
			logger.Warn("absent API key")
			http.Error(res, "Absent API key", http.StatusUnauthorized)
			return
		}

		if !validKeys[key] {
			logger.WithField("api_key", key).Warn("wrong API key")
			http.Error(res, "Wrong API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}
