package internalhttp

import (
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	BytesLength int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.BytesLength += n

	return n, err
}

func loggingMiddleware(next http.Handler, log app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myWriter := &ResponseWriter{w, 0, 0}
		next.ServeHTTP(myWriter, r)
		log.LogHTTPRequest(r, myWriter.StatusCode, myWriter.BytesLength)
	})
}
