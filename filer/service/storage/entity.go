package storage

import (
	"context"
	"github.com/gocastsian/roham/types"
	"io"
	"time"
)

type Provider interface {
	GetFileContent(ctx context.Context, bucketName, key string) (io.ReadCloser, error)
	GeneratePreSignedURL(bucketName, key string, duration time.Duration) (string, error)
	CreateBucket(ctx context.Context, bucketName string) error
}

type FileMetadata struct {
	ID         types.ID          `json:"id"`
	Key        string            `json:"key"`
	BucketName string            `json:"bucket_name"`
	Metadata   map[string]string `json:"metadata"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type Bucket struct {
	ID   types.ID `json:"id"`
	Name string   `json:"bucket_name"`
}

type Storage struct {
	Category string            `json:"category"` // avatar | layers
	Type     string            `json:"type"`     // filesystem | s3
	Name     string            `json:"name"`
	Config   map[string]string `json:"config"`
}
