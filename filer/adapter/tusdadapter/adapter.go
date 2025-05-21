package tusdadapter

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gocastsian/roham/filer/storageprovider"
	"github.com/gocastsian/roham/filer/storageprovider/s3storage"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
	"log"
	"net/http"
	"os"
	"regexp"
)

type UploadValidator interface {
	ValidateUpload(targetBucket string, size int64) error
}

func NewHandlerWithS3Store(storageName string, s3 *s3.S3, v UploadValidator) (*tusd.Handler, error) {

	tusd.NewStoreComposer()
	store := s3store.New(storageName, s3)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:      "/uploads/",
		StoreComposer: composer,
		PreUploadCreateCallback: func(hook tusd.HookEvent) error {
			targetBucket := hook.HTTPRequest.Header.Get("Bucket")
			if targetBucket == "" {
				return errors.New("bucket name is required")
			}
			return v.ValidateUpload(targetBucket, hook.Upload.Size)
		},
		NotifyCompleteUploads:   true,
		NotifyTerminatedUploads: true,
		NotifyUploadProgress:    true,
		NotifyCreatedUploads:    true,
	})

	if err != nil {
		return nil, err
	}

	return handler, err
}

func NewHandlerWithFileStore(storageName, basePath string, v UploadValidator) (*tusd.Handler, error) {

	uploadsDir := fmt.Sprintf("%s/%s", basePath, storageName)
	store := filestore.New(uploadsDir)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	//todo get cors from configs. cors should be handle by proxy or tusd handle?

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath: "/uploads/",
		Cors: &tusd.CorsConfig{
			AllowOrigin:   regexp.MustCompile(".*"),
			AllowHeaders:  "*",
			ExposeHeaders: "*",
			AllowMethods:  "*",
		},
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		PreUploadCreateCallback: func(hook tusd.HookEvent) error {

			targetBucket := hook.HTTPRequest.Header.Get("Bucket")
			if targetBucket == "1" {
				return tusd.NewHTTPError(errors.New("missing required header: Bucket"), http.StatusUnprocessableEntity)
			}
			return v.ValidateUpload(targetBucket, hook.Upload.Size)
		},
		NotifyTerminatedUploads: true,
		NotifyUploadProgress:    true,
		NotifyCreatedUploads:    true,
	})

	// Create uploads dir if not exist
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadsDir, os.ModePerm); err != nil {
			log.Fatalf("Could not create upload dir: %v", err)
		}
	}

	if err != nil {
		return nil, err
	}

	return handler, err
}

func New(p storageprovider.Provider, v UploadValidator) (*tusd.Handler, error) {

	switch p.Config().Type {
	case "filesystem":
		return NewHandlerWithFileStore(p.Config().TempStorage, p.Config().BasePath, v)
	case "s3":
		store, ok := p.(*s3storage.Storage)
		if !ok {
			return nil, errors.New("not a s3 storage")
		}
		return NewHandlerWithS3Store(p.Config().TempStorage, store.S3(), v)
	}

	return nil, errors.New("unknown storage type")
}
