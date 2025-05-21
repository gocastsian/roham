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
		cfg: cfg,
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

	url := fmt.Sprintf("http://localhost:5006/storages/%s/files/%s", storageName, key)
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

func (s *Storage) Config() storageprovider.StorageConfig {
	return s.cfg
}
