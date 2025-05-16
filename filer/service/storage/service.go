package storage

import (
	"context"
	"github.com/gocastsian/roham/types"
	"io"
	"time"
)

type Service struct {
	provider    Provider
	storageRepo StorageRepo
	fileRepo    FileMetadataRepo
}

type StorageRepo interface {
	Insert(ctx context.Context, s CreateStorageInput) (types.ID, error)
}

type FileMetadataRepo interface {
	Insert(ctx context.Context, fileMetadata CreateFileMetadataInput) (types.ID, error)
}

func NewStorageService(p Provider, fr FileMetadataRepo, br StorageRepo) Service {
	return Service{
		provider:    p,
		fileRepo:    fr,
		storageRepo: br,
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

func (s Service) CreateStorage(ctx context.Context, input CreateStorageInput) (*CreateStorageOutput, error) {

	err := s.provider.MakeStorage(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	id, err := s.storageRepo.Insert(ctx, input)
	if err != nil {
		return nil, err
	}

	return &CreateStorageOutput{
		ID:   id,
		Name: input.Name,
		Kind: input.Kind,
	}, nil
}
