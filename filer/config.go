package filer

import (
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/logger"
	"time"
)

type Config struct {
	HTTPServer           httpserver.Config `koanf:"http_server"`
	Logger               logger.Config     `koanf:"logger"`
	TotalShutdownTimeout time.Duration     `koanf:"total_shutdown_timeout"`
	MinioStorage         MinioStorage      `koanf:"minio_storage"`
	Uploader             Uploader          `koanf:"uploader"`
}

type Uploader struct {
	HTTPServer httpserver.Config `koanf:"http_server"`
	Logger     logger.Config     `koanf:"logger"`
}

type MinioStorage struct {
	Endpoint  string `koanf:"endpoint"`
	AccessKey string `koanf:"access_key"`
	SecretKey string `koanf:"secret_key"`
}
