package server

import (
	"context"
	"net/http"
	"time"
)

type ServerConfig struct {
	Host string
	Port string
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg ServerConfig, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           cfg.Host + ":" + cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
