package storage

import (
	"context"
	"github.com/gocastsian/roham/filer/storageprovider"
	"github.com/gocastsian/roham/types"
	"io"
	"time"
)

type Service struct {
	provider    storageprovider.Provider
	storageRepo Repository
	fileRepo    FileMetadataRepo
}

type Repository interface {
	Insert(ctx context.Context, s CreateStorageInput) (types.ID, error)
}

type FileMetadataRepo interface {
	InsertFileMetadata(ctx context.Context, fileMetadata FileMetadata) (types.ID, error)
}

func NewStorageService(p storageprovider.Provider, fr FileMetadataRepo, r Repository) Service {
	return Service{
		provider:    p,
		fileRepo:    fr,
		storageRepo: r,
	}
}

func (s Service) GetFile(ctx context.Context, storageName, fileKey string) (io.ReadCloser, error) {
	return s.provider.GetFile(ctx, storageName, fileKey)
}

func (s Service) GeneratePreSignedURL(ctx context.Context, storageName, fileKey string, t time.Duration) (string, error) {
	return s.provider.GeneratePreSignedURL(storageName, fileKey, t)
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
