package userapp

import (
	"time"

	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/logger"
	"github.com/gocastsian/roham/pkg/postgresql"
	"github.com/gocastsian/roham/userapp/repository"
	"github.com/gocastsian/roham/userapp/service/guard"
)

type Config struct {
	HTTPServer           httpserver.Config `koanf:"http_server"`
	PostgresDB           postgresql.Config `koanf:"postgres_db"`
	Repository           repository.Config `koanf:"repository"`
	Logger               logger.Config     `koanf:"logger"`
	TotalShutdownTimeout time.Duration     `koanf:"total_shutdown_timeout"`
	Guard                guard.Config      `koanf:"guard"`
}
