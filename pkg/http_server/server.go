package httpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port               int           `koanf:"port"`
	Cors               Cors          `koanf:"cors"`
	ShutDownCtxTimeout time.Duration `koanf:"shutdown_context_timeout"`
}

type Cors struct {
	AllowOrigins []string `koanf:"allow_origins"`
	AllowHeaders []string `koanf:"allow_headers"`
	Skip         bool     `koanf:"skip"`
}

type Server struct {
	Router *echo.Echo
	Config Config
}

func New(cfg Config) Server {
	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper: func(c echo.Context) bool {
			return cfg.Cors.Skip
		},
		AllowOrigins: cfg.Cors.AllowOrigins,
		//AllowHeaders: cfg.Cors.AllowHeaders,
	}))

	return Server{
		Router: e,
		Config: cfg,
	}
}

// register custom handler
func (s Server) RegisterHandler(route string, handler echo.HandlerFunc) {
	s.Router.GET(route, handler)
}

// start server
func (s Server) Start() error {
	return s.Router.Start(fmt.Sprintf(":%d", s.Config.Port))
}

func (s Server) Stop(ctx context.Context) error {
	return s.Router.Shutdown(ctx)
}
