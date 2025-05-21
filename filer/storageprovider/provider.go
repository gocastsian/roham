package storageprovider

import (
	"context"
	"io"
	"time"
)

type Provider interface {
	GetFile(ctx context.Context, storageName, fileKey string) (io.ReadCloser, error)
	GeneratePreSignedURL(storageName, fileKey string, duration time.Duration) (string, error)
	MakeStorage(ctx context.Context, name string) error
	Config() StorageConfig
}

// StorageConfig holds configuration for file StorageConfig, supporting both local filesystem and S3.
type StorageConfig struct {
	// Type defines the StorageConfig driver ("filesystem" or "s3").
	Type string `koanf:"type"`

	// TempStorage is the default bucket/path where files are initially uploaded.
	// After processing, files may be moved to a permanent location.
	TempStorage string `koanf:"temp_storage"`

	// --- Filesystem-specific settings ---
	BasePath string `koanf:"base_path"`

	// --- S3-specific settings ---
	Region    string `koanf:"region"`
	Endpoint  string `koanf:"endpoint"`
	AccessKey string `koanf:"access_key"`
	SecretKey string `koanf:"secret_key"`
}
