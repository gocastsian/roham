package storage

import (
	"context"
	"github.com/gocastsian/roham/types"
	"io"
	"time"
)

type Service struct {
	provider   Provider
	bucketRepo BucketRepo
	fileRepo   FileMetadataRepo
}

type BucketRepo interface {
	CreateBucket(ctx context.Context, name string) (types.ID, error)
}

type FileMetadataRepo interface {
	Create(ctx context.Context, fileMetadata FileMetadata) (types.ID, error)
}

func NewStorageService(p Provider, fr FileMetadataRepo, br BucketRepo) Service {
	return Service{
		provider:   p,
		fileRepo:   fr,
		bucketRepo: br,
	}
}

func (s Service) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	// todo get bucket name using FileRepo
	// todo send metadata of file in addition to content
	bucketName := "default-bucket"
	return s.provider.GetFileContent(ctx, bucketName, key)
}

func (s Service) GeneratePreSignedURL(ctx context.Context, key string, t time.Duration) (string, error) {
	// todo get bucket name using FileRepo
	bucketName := "default-bucket"
	return s.provider.GeneratePreSignedURL(bucketName, key, t)
}

func (s Service) CreateBucket(ctx context.Context, name string) (*Bucket, error) {

	err := s.provider.CreateBucket(ctx, name)
	if err != nil {
		return nil, err
	}
	bucketID, err := s.bucketRepo.CreateBucket(ctx, name)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		ID:   bucketID,
		Name: name,
	}, nil
}
