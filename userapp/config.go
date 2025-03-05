package userapp

import (
	"time"

	httpserver "roham/pkg/http_server"
	"roham/pkg/logger"
	"roham/pkg/postgresql"
	"roham/userapp/repository"
)

type Config struct {
	HTTPServer           httpserver.Config `koanf:"http_server"`
	PostgresDB           postgresql.Config `koanf:"postgres_db"`
	Repository           repository.Config `koanf:"repository"`
	Logger               logger.Config     `koanf:"logger"`
	TotalShutdownTimeout time.Duration     `koanf:"total_shutdown_timeout"`
}
