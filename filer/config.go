package filer

import (
	"github.com/gocastsian/roham/filer/storageprovider"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/logger"
	"github.com/gocastsian/roham/pkg/postgresql"
	"time"
)

type Config struct {
	HTTPServer           httpserver.Config             `koanf:"http_server"`
	Logger               logger.Config                 `koanf:"logger"`
	PostgresDB           postgresql.Config             `koanf:"postgres_db"`
	TotalShutdownTimeout time.Duration                 `koanf:"total_shutdown_timeout"`
	Uploader             uploader                      `koanf:"uploader"`
	Storage              storageprovider.StorageConfig `koanf:"storage"`
}

type uploader struct {
	HTTPServer httpserver.Config `koanf:"server"`
	Logger     logger.Config     `koanf:"logger"`
}
