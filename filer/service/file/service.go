package file

import (
	"context"
	"io"
	"time"
)

type Service struct {
	storage Storage
}

type Storage interface {
	GetFileContent(ctx context.Context, bucketName, key string) (io.ReadCloser, error)
	GeneratePreSignedURL(bucketName, key string, duration time.Duration) (string, error)
}

func NewFileService(s Storage) Service {
	return Service{
		storage: s,
	}
}

func (s Service) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	// todo get bucket name using FileRepo
	// todo send metadata of file in addition to content
	bucketName := "default-bucket"
	return s.storage.GetFileContent(ctx, bucketName, key)
}

func (s Service) GeneratePreSignedURL(ctx context.Context, key string, t time.Duration) (string, error) {
	// todo get bucket name using FileRepo
	bucketName := "default-bucket"
	return s.storage.GeneratePreSignedURL(bucketName, key, t)
}
