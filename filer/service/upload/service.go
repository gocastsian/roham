package upload

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gocastsian/roham/filer/service/filestorage"
	"github.com/gocastsian/roham/filer/storageprovider"
	"github.com/gocastsian/roham/types"
)

type Service struct {
	logger           *slog.Logger
	fileMetadataRepo FileMetadataRepo
	storageFinder    StorageFinder
	storageProvider  storageprovider.Provider
}

func NewUploadService(l *slog.Logger, sp storageprovider.Provider, fileMetadataRepo FileMetadataRepo, storageRepo StorageFinder) Service {
	return Service{
		logger:           l,
		fileMetadataRepo: fileMetadataRepo,
		storageFinder:    storageRepo,
		storageProvider:  sp,
	}
}

type StorageFinder interface {
	FindByID(ctx context.Context, id types.ID) (filestorage.Storage, error)
}

type FileMetadataRepo interface {
	InsertFileMetadata(ctx context.Context, fileMetadata filestorage.FileMetadata) (types.ID, error)
}

func (s *Service) OnCompletedUploads(ctx context.Context, i filestorage.CreateFileMetadataInput) error {

	storage, err := s.storageFinder.FindByID(ctx, i.TargetStorageID)
	if err != nil {
		return err
	}

	err = s.storageProvider.MoveFileToStorage(i.FileKey, s.storageProvider.Config().TempStorage, storage.Name)
	if err != nil {
		return err
	}

	newFileMetadata := filestorage.FileMetadata{
		StorageID: i.TargetStorageID,
		FileKey:   i.FileKey,
		FileName:  i.FileName,
		MimeType:  i.MimeType,
		FileSize:  i.Size,
		CreatedAt: time.Time{},
		UpdatedAt: nil,
	}
	_, err = s.fileMetadataRepo.InsertFileMetadata(ctx, newFileMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateUpload(targetStorageID types.ID, mimeType string, size int64) error {

	ctx := context.Background()
	storage, err := s.storageFinder.FindByID(ctx, targetStorageID)
	if err != nil {
		return err
	}
	const (
		GB = 1024 * 1024 * 1024
		MB = 1024 * 1024
		KB = 1024
	)
	// todo get from validation config
	uploadConstraints := map[string]struct {
		MaxSize      int64
		AllowedTypes []string
	}{
		"avatar": {
			MaxSize:      5 * MB,
			AllowedTypes: []string{"image/jpeg", "image/jpg"},
		},
		"map-layer": {
			MaxSize:      2 * GB,
			AllowedTypes: []string{}, // empty means all types allowed or define specific ones
		},
	}

	// Check if storage kind is valid
	constraint, exists := uploadConstraints[storage.Kind]
	if !exists {
		return fmt.Errorf("unsupported storage kind: %s", storage.Kind)
	}

	// Validate file size
	if size > constraint.MaxSize {
		return fmt.Errorf("file size exceeds maximum allowed (%d bytes)", constraint.MaxSize)
	}

	// Validate MIME type if constraints are defined
	if len(constraint.AllowedTypes) > 0 {
		validType := false
		for _, allowedType := range constraint.AllowedTypes {
			if mimeType == allowedType {
				validType = true
				break
			}
		}
		if !validType {
			return fmt.Errorf("unsupported file type: %s. Allowed types: %v",
				mimeType, constraint.AllowedTypes)
		}
	}

	return nil
}
