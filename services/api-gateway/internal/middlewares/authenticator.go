package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func AuthenticatorMiddleware(h http.Handler, l *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			h.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = r.URL.Query().Get("api_key")
		}

		validKeys := map[string]bool{
			"secret-key-1234": true,
			"test-key-1234":   true,
			"dev-key-1234":    true,
		}

		if apiKey == "" {
			l.Warn("Absent API key")
			http.Error(w, "Provide API key or go fuck yourself", http.StatusUnauthorized)
			return
		}

		if !validKeys[apiKey] {
			l.WithField("key", apiKey).Warn("Wrong API key")
			http.Error(w, "Wrong API key", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
