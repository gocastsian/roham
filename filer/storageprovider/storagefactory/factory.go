package storagefactory

import (
	"errors"
	"fmt"
	"github.com/gocastsian/roham/filer/storageprovider"
	"github.com/gocastsian/roham/filer/storageprovider/filestorage"
	"github.com/gocastsian/roham/filer/storageprovider/s3storage"
)

func New(cfg storageprovider.StorageConfig) (storageprovider.Provider, error) {

	switch cfg.Type {
	case "filesystem":
		return filestorage.New(cfg)
	case "s3":
		return s3storage.New(cfg)
	}

	return nil, errors.New(fmt.Sprintf("Invalid storage type: %s", cfg.Type))
}
