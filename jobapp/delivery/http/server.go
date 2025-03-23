package http

import (
	"context"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
)

type Server struct {
	HTTPServer httpserver.Server
}

func New(server httpserver.Server) Server {
	return Server{
		HTTPServer: server,
	}
}

func (s Server) Server() error {
	s.RegisterRoutes()

	if err := s.HTTPServer.Start(); err != nil {
		return err
	}
	return nil
}

func (s Server) Stop(ctx context.Context) error {
	return s.HTTPServer.Stop(ctx)
}

func (s Server) RegisterRoutes() {
	v1 := s.HTTPServer.Router.Group("/v1")

	v1.GET("/health-check", s.healthCheck)
}
