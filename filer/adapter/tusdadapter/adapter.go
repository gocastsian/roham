package tusdadapter

import (
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
)

type PreUploadCreateCallback func(hook tusd.HookEvent) error

type UploadEvent struct {
}

type Adapter struct {
}

func NewWithFileStore(store filestore.FileStore, callback PreUploadCreateCallback) (*tusd.Handler, error) {

	tusd.NewStoreComposer()

	// Create a new tusd handler with the file store
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:                "/uploads/",
		StoreComposer:           composer,
		NotifyCompleteUploads:   true,
		PreUploadCreateCallback: callback,
	})

	if err != nil {
		return nil, err
	}

	return handler, err
}

func NewWithS3Store(store s3store.S3Store, callback PreUploadCreateCallback) (*tusd.Handler, error) {

	tusd.NewStoreComposer()

	// Create a new tusd handler with the file store
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:                "/uploads/",
		StoreComposer:           composer,
		NotifyCompleteUploads:   true,
		PreUploadCreateCallback: callback,
	})

	if err != nil {
		return nil, err
	}

	return handler, err
}
