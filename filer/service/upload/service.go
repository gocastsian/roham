package upload

import (
	"context"
	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/types"
	"log/slog"
	"strconv"
	"time"
)

type Service struct {
	logger           *slog.Logger
	fileMetadataRepo FileMetadataRepository
	fileMover        fileMover
}

type fileMover interface {
	moveTo(fileKey, to string) error
}

func NewUploadService(l *slog.Logger, fileMetadataRepo FileMetadataRepository) Service {
	return Service{
		logger:           l,
		fileMetadataRepo: fileMetadataRepo,
	}
}

type FileMetadataRepository interface {
	Create(ctx context.Context, fileMetadata storage.FileMetadata) (types.ID, error)
}

func (s *Service) OnCompletedUploads(ctx context.Context, uploadID, bucketName string, metaData map[string]string) error {

	s.logger.Info("Upload completed :" + uploadID)

	f := storage.FileMetadata{
		Key:        uploadID,
		BucketName: bucketName,
		Metadata:   metaData,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	_, err := s.fileMetadataRepo.Create(ctx, f)
	if err != nil {
		return err
	}

	//todo move file to its original bucket.

	return nil
}

func (s *Service) ValidateUpload(targetBucket string, size int64) error {
	s.logger.Info("Validate upload request. size is : " + strconv.Itoa(int(size)))

	//todo validate base on

	return nil
}
