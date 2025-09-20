package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

func LogrusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use a response writer that captures status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		level := log.GetLevel()

		if ww.Status() >= 500 {
			level = log.ErrorLevel
		} else if ww.Status() >= 400 {
			level = log.WarnLevel
		}

		log.WithFields(log.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      ww.Status(),
			"duration_ms": duration.Milliseconds(),
			"remote_addr": r.RemoteAddr,
		}).Log(level, "incoming request")
	})
}
