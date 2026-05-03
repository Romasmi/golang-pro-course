package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/server/http/handlers"
)

type Server struct {
	httpServer *http.Server
	logger     Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface{}

func NewServer(logger Logger, _ Application, host, port string) *Server {
	s := &Server{
		logger: logger,
	}

	mux := http.NewServeMux()
	s.registerHandlers(mux)

	handler := loggingMiddleware(logger, mux)

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

func (s *Server) registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/hello", handlers.HelloHandler)
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info("starting http server on " + s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server failed: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping http server")
	return s.httpServer.Shutdown(ctx)
}
