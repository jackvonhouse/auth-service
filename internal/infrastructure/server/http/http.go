package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackvonhouse/auth-service/config"
)

type Server struct {
	server *http.Server
}

func New(
	handler http.Handler,
	config *config.ServerHTTP,
) *Server {

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: handler,
	}

	return &Server{
		server: &httpServer,
	}
}

func (s *Server) Run() error {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(
	ctx context.Context,
) error {

	return s.server.Shutdown(ctx)
}
