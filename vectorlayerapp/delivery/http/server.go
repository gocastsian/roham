package http

import (
	"context"
	echomiddleware "github.com/gocastsian/roham/pkg/echo_middleware"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/types"
	"log/slog"
)

type Server struct {
	HTTPServer   httpserver.Server
	tileHandler  GenLayerHandler
	layerHandler ImportLayerHandler
	logger       *slog.Logger
}

func New(server httpserver.Server, tileHandler GenLayerHandler, layerHandler ImportLayerHandler, logger *slog.Logger) Server {
	return Server{
		HTTPServer:   server,
		tileHandler:  tileHandler,
		layerHandler: layerHandler,
		logger:       logger,
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

	tileServiceGroup := v1.Group("/tile")
	{
		tileServiceGroup.GET("/health-check", s.tileHandler.healthCheck)
	}

	layerServiceGroup := v1.Group("/layer")
	{
		layerServiceGroup.GET("/health-check", s.layerHandler.healthCheck)
		layerServiceGroup.POST("/create-job", s.layerHandler.createJob, echomiddleware.ParseUserDataMiddleware, echomiddleware.AccessCheck([]types.Role{types.RoleAdmin}))
	}
}
