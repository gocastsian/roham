package jobapp

import (
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/logger"
	"time"
)

type Config struct {
	HTTPServer           httpserver.Config `koanf:"http_server"`
	Logger               logger.Config     `koanf:"logger"`
	TotalShutdownTimeout time.Duration     `koanf:"total_shutdown_timeout"`
}
