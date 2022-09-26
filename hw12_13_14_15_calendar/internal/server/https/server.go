package https

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	host   string
	port   string
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(message string, params ...interface{})
	Error(message string, params ...interface{})
	LogRequest(r *http.Request, code, length int)
}

type Application interface{}

func NewServer(logger Logger, app *app.App, host, port string) *Server {
	httpServer := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	newServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: loggingMiddleware(routes(app), logger),
	}

	httpServer.server = newServer

	return httpServer
}

func routes(app *app.App) http.Handler {
	r := mux.NewRouter()
	eventHandler := NewEventHandler(app)

	r.HandleFunc("/events", eventHandler.List).Methods("GET")
	r.HandleFunc("/events", eventHandler.Create).Methods("POST")
	r.HandleFunc("/events/{id}", eventHandler.Update).Methods("PUT")
	r.HandleFunc("/events/{id}", eventHandler.Delete).Methods("DELETE")

	return r
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("HTTP server run %s:%s", s.host, s.port)

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server stopped")

	return s.server.Shutdown(ctx)
}
