package filestorage

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/gocastsian/roham/filer/storageprovider"
	"github.com/gocastsian/roham/types"
)

type Service struct {
	logger          *slog.Logger
	storageProvider storageprovider.Provider
	storageRepo     StorageRepository
	fileRepo        FileMetadataRepo
}

type StorageRepository interface {
	Insert(ctx context.Context, s CreateStorageInput) (types.ID, error)
	FindByID(ctx context.Context, id types.ID) (Storage, error)
}

type FileMetadataRepo interface {
	InsertFileMetadata(ctx context.Context, fileMetadata FileMetadata) (types.ID, error)
	FindByKey(ctx context.Context, key string) (FileMetadata, error)
}

func NewStorageService(l *slog.Logger, p storageprovider.Provider, fr FileMetadataRepo, r StorageRepository) Service {
	return Service{
		logger:          l,
		storageProvider: p,
		fileRepo:        fr,
		storageRepo:     r,
	}
}

func (s Service) GetFileByKey(ctx context.Context, fileKey string) (io.ReadCloser, error) {

	fileMetadata, err := s.fileRepo.FindByKey(ctx, fileKey)
	if err != nil {
		return nil, err
	}

	var storageName string

	storage, err := s.storageRepo.FindByID(ctx, fileMetadata.StorageID)
	if err != nil {
		return nil, err
	}
	storageName = storage.Name

	return s.storageProvider.GetFile(ctx, storageName, fileKey)
}

func (s Service) GeneratePreSignedURL(ctx context.Context, storageName, fileKey string, t time.Duration) (string, error) {
	return s.storageProvider.GeneratePreSignedURL(storageName, fileKey, t)
}

func (s Service) CreateStorage(ctx context.Context, input CreateStorageInput) (*CreateStorageOutput, error) {

	err := s.storageProvider.MakeStorage(ctx, input.Name)
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
