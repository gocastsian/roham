package vectorlayerapp

import (
	"github.com/gocastsian/roham/adapter/temporal"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/logger"
	"github.com/gocastsian/roham/pkg/postgresql"
	"time"
)

type Config struct {
	HTTPServer           httpserver.Config `koanf:"http_server"`
	PostgresDB           postgresql.Config `koanf:"postgres_db"`
	Logger               logger.Config     `koanf:"logger"`
	TotalShutdownTimeout time.Duration     `koanf:"total_shutdown_timeout"`
	Temporal             temporal.Config
}
