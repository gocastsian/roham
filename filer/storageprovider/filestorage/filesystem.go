package filestorage

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/filer/storageprovider"
	"io"
	"os"
	"time"
)

type Storage struct {
	basePath string
	cfg      storageprovider.StorageConfig
}

func New(cfg storageprovider.StorageConfig) (*Storage, error) {
	return &Storage{
		cfg:      cfg,
		basePath: cfg.BasePath,
	}, nil
}

func (s *Storage) GetFile(ctx context.Context, storageName, fileKey string) (io.ReadCloser, error) {

	filePath := fmt.Sprintf("%s/%s/%s", s.basePath, storageName, fileKey)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *Storage) GeneratePreSignedURL(storageName, key string, duration time.Duration) (string, error) {

	//todo implement custom pre-signed url using database
	url := fmt.Sprintf("http://localhost:5006/v1/files/%s/download", storageName+":"+key)

	return url, nil

}

func (s *Storage) MakeStorage(ctx context.Context, name string) error {

	newFolder := fmt.Sprintf("%s/%s", s.basePath, name)
	err := os.MkdirAll(newFolder, 0755)
	if err != nil {

		return err
	}

	return nil
}

func (s *Storage) MoveFileToStorage(fileKey, fromStorageName, toStorageName string) error {

	// Define the destination directory
	destDir := fmt.Sprintf("%s/%s", s.basePath, toStorageName)

	// Ensure the destination directory exists
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	sourcePath := fmt.Sprintf("%s/%s/%s", s.basePath, fromStorageName, fileKey)

	destPath := fmt.Sprintf("%s/%s/%s", s.basePath, toStorageName, fileKey)

	// Move the file
	err = os.Rename(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	err = os.Rename(sourcePath+".info", destPath+".info")
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	return nil
}

func (s *Storage) Config() storageprovider.StorageConfig {
	return s.cfg
}
