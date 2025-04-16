package http

import (
	"context"
	"log/slog"

	httpserver "github.com/gocastsian/roham/pkg/http_server"
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
	v1.GET("/health-check", s.healthCheck)
	v1.GET("/auth", s.Handler.authenticate)
	v1.GET("/authz", s.Handler.authorize)

	userGroup := v1.Group("/users")
	userGroup.GET("/", s.Handler.GetAllUsers)
	userGroup.POST("/login", s.Handler.Login)
	userGroup.POST("", s.Handler.registerUser)
}
