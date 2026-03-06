package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(h http.Handler, l *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received := time.Now()

		l.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		}).Info("Request received")

		h.ServeHTTP(w, r)

		processed := time.Since(received)

		l.WithFields(logrus.Fields{
			"duration": processed,
			"path":     r.URL.Path,
		}).Info("Request processed")
	})
}
