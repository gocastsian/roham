package upload

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/types"
	"log/slog"
	"strconv"
)

type Service struct {
	logger           *slog.Logger
	fileMetadataRepo FileMetadataRepo
	fileMover        fileMover
}

type fileMover interface {
	moveTo(fileKey, to string) error
}

func NewUploadService(l *slog.Logger, fileMetadataRepo FileMetadataRepo) Service {
	return Service{
		logger:           l,
		fileMetadataRepo: fileMetadataRepo,
	}
}

type FileMetadataRepo interface {
	Insert(ctx context.Context, fileMetadata storage.CreateFileMetadataInput) (types.ID, error)
}

func (s *Service) OnCompletedUploads(ctx context.Context, input storage.CreateFileMetadataInput) error {

	s.logger.Info("Upload completed. FileName : " + input.FileName)

	fmt.Println(input)
	_, err := s.fileMetadataRepo.Insert(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateUpload(targetBucket string, size int64) error {
	s.logger.Info("Validate upload request. size is : " + strconv.Itoa(int(size)))

	//todo validate base on

	return nil
}
