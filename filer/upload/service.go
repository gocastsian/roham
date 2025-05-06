package upload

import (
	"log/slog"
	"strconv"
)

type Service struct {
	logger *slog.Logger
}

func NewUploadService(l *slog.Logger) Service {
	return Service{
		logger: l,
	}
}

func (s *Service) CompleteUpload(uploadID string, metaData map[string]string) error {

	s.logger.Info("Upload completed :" + uploadID)
	//todo save metadata to file db
	//todo publish event to temporal

	return nil
}

func (s *Service) ValidateUpload(UploadRequest NewUploadRequest) error {
	s.logger.Info("Validate upload request. size is : " + strconv.Itoa(int(UploadRequest.Size)))

	//todo validate base on

	return nil
}
