package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	internalapp "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	host   string
	port   string
	logger internalapp.Logger
	server *http.Server
}

func NewServer(logger internalapp.Logger, app *internalapp.App, host, port string) *Server {
	server := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: loggingMiddleware(NewRouter(app), logger),
	}

	server.server = httpServer

	return server
}

func NewRouter(app *internalapp.App) http.Handler {
	handlers := NewServerHandlers(app)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HelloWorld).Methods("GET")
	r.HandleFunc("/events", handlers.CreateEvent).Methods("POST")
	r.HandleFunc("/events/{id}", handlers.UpdateEvent).Methods("PUT")
	r.HandleFunc("/events/{id}", handlers.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events", handlers.ListEvents).Methods("GET")

	return r
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("HTTP server listen and serve %s:%s", s.host, s.port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
