package upload

import (
	"context"
	"errors"
	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/types"
	"log/slog"
	"strconv"
	"time"
)

type Service struct {
	logger           *slog.Logger
	fileMetadataRepo FileMetadataRepo
	storageFinder    StorageFinder
}

func NewUploadService(l *slog.Logger, fileMetadataRepo FileMetadataRepo, storageRepo StorageFinder) Service {
	return Service{
		logger:           l,
		fileMetadataRepo: fileMetadataRepo,
		storageFinder:    storageRepo,
	}
}

type StorageFinder interface {
	FindByName(ctx context.Context, name string) (storage.Storage, error)
}

type FileMetadataRepo interface {
	InsertFileMetadata(ctx context.Context, fileMetadata storage.FileMetadata) (types.ID, error)
}

func (s *Service) OnCompletedUploads(ctx context.Context, input storage.CreateFileMetadataInput) error {

	//s.logger.Info("Upload completed. FileName : " + input.FileName)

	targetStorage, err := s.storageFinder.FindByName(ctx, input.TargetStorageName)

	if err != nil {
		return err
	}

	newFileMetadata := storage.FileMetadata{
		StorageID: targetStorage.ID,
		FileKey:   "",
		FileName:  "",
		MimeType:  "",
		Size:      "",
		CreatedAt: time.Time{},
		ClaimedAt: nil,
	}
	_, err = s.fileMetadataRepo.InsertFileMetadata(ctx, newFileMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateUpload(storageKind string, size int64) error {
	s.logger.Info("Validate upload request. size is : " + strconv.Itoa(int(size)))

	//todo get storage kinds from config
	switch storageKind {
	case "avatar":
		if size > 10240 {
			return errors.New("file size is too big")
		}
	case "map-layer":
		if size > 1024000 {
			return errors.New("file size is too big")
		}
	}

	return nil
}

func (s *Service) OnFileClaimed(ctx context.Context, fileKey string) error {

	//todo move file to target storage
	//todo update claimed_at of fileMetadata
	return nil
}
