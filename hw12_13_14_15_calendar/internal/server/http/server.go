package http

import (
	"context"
	"net"
	netHttp "net/http"
)

type Server struct {
	host   string
	port   string
	logger Logger
	server *netHttp.Server
}

type Logger interface {
	Info(message string, params ...interface{})
	Error(message string, params ...interface{})
	LogRequest(r *netHttp.Request, code, length int)
}

type Application interface {
}

func NewServer(logger Logger, app Application, host, port string) *Server {
	httpServer := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	newServer := &netHttp.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: loggingMiddleware(netHttp.HandlerFunc(httpServer.handleHTTP), logger),
	}

	httpServer.server = newServer

	return httpServer
}

func (s *Server) handleHTTP(w netHttp.ResponseWriter, r *netHttp.Request) {
	w.Write([]byte("hello-world"))
	w.WriteHeader(200)
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
	return s.server.Shutdown(ctx)
}
