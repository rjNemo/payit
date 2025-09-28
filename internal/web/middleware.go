package web

import (
	"log"
	"net/http"
	"time"
)

type WrappedWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &WrappedWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.StatusCode, time.Since(start))
	})
}
