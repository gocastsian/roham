package filer

import (
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/logger"
	"github.com/gocastsian/roham/pkg/postgresql"
	"time"
)

type Config struct {
	HTTPServer           httpserver.Config `koanf:"http_server"`
	Logger               logger.Config     `koanf:"logger"`
	PostgresDB           postgresql.Config `koanf:"postgres_db"`
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

type StorageConfig struct {
	BucketName      string `json:"bucket_name"`
	StorageProvider string `json:"storage_provider"`
	// filesystem
	BasePath string `json:"base_path,omitempty"`
	//s3
	Region    string `json:"region,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
}
