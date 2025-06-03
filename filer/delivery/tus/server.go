package tus

import (
	"context"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/labstack/echo/v4"
	tusd "github.com/tus/tusd/pkg/handler"
	"log/slog"
	"net/http"
)

type Server struct {
	Handler    Handler
	HTTPServer httpserver.Server
	logger     *slog.Logger
}

func NewServer(l *slog.Logger, httpServer httpserver.Server, h Handler) Server {
	return Server{
		logger:     l,
		HTTPServer: httpServer,
		Handler:    h,
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
	r := s.HTTPServer.Router.Group("")
	r.GET("/health-check", s.healthCheck)
	r.Any("/uploads/*", echo.WrapHandler(http.StripPrefix("/uploads/", s.Handler.TusHandler)))
}

func (s Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "everything is good!",
	})
}

func (s Server) EventHandler() *tusd.Handler {
	return s.Handler.TusHandler
}
