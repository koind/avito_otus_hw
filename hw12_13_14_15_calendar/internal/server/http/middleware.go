package http

import (
	netHttp "net/http"
)

type ResponseWriter struct {
	netHttp.ResponseWriter
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

func loggingMiddleware(next netHttp.Handler, log Logger) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		writer := &ResponseWriter{w, 0, 0}
		next.ServeHTTP(writer, r)
		log.LogRequest(r, writer.StatusCode, writer.BytesLength)
	})
}
