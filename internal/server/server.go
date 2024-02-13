package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
)

// Server represents a http server structure.
type Server struct {
	httpServer *http.Server
}

// NewServer is a Server constructor.
func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderBytes << 20,
		},
	}
}

// Run runs a http server.
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Stop stops a http server.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
