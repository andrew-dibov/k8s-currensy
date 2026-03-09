package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log := logger.WithFields(logrus.Fields{
			"req_method": req.Method,
			"req_path":   req.URL.Path,
			"req_addr":   req.RemoteAddr,
		})

		received := time.Now()
		log.WithField("req_received", received).Info("request received")

		next.ServeHTTP(res, req)

		duration := time.Since(received)
		log.WithField("req_duration", duration).Info("request processed")
	})
}
