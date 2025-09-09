package middleware_logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct{
	http.ResponseWriter 
	statusCode int
}


func ZapLogger(logger *zap.Logger) func( next http.Handler) http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)

			latency := time.Since(start)
			logger.Info("incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.statusCode),
				zap.Duration("latency", latency),
				zap.String("client_ip", r.RemoteAddr),
			)
		})
	}
}

func(rw *responseWriter) WriteHeader(code int){
	rw.statusCode= code
	rw.ResponseWriter.WriteHeader(code)
}