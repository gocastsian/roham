package http

import (
	"context"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"log/slog"
)

type Server struct {
	HTTPServer httpserver.Server
	Handler    Handler
	logger     *slog.Logger
}

func New(server httpserver.Server, handler Handler, logger *slog.Logger) Server {
	return Server{
		HTTPServer: server,
		Handler:    handler,
		logger:     logger,
	}
}

func (s Server) Serve() error {
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
	v1.GET("/health-check", s.Handler.healthCheck)

	layerGroup := v1.Group("/layer")
	layerGroup.GET("/import", s.Handler.ImportLayer)
}
