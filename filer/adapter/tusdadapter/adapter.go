package tusdadapter

import (
	"errors"
	"github.com/gocastsian/roham/filer/adapter/s3adapter"
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

func New() *tusd.Handler {
	// Set up tusd file store
	store := filestore.FileStore{
		Path: "./uploads/",
	}

	// Create tusd handler config
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	h, err := tusd.NewHandler(tusd.Config{
		BasePath:                "/files/",
		StoreComposer:           composer,
		RespectForwardedHeaders: true,
	})
	if err != nil {
		log.Fatalf("Failed to create tusd handler: %v", err)
	}
	return h
}

func NewHandlerWithS3Store(bucketName string, v UploadValidator, adapter *s3adapter.Adapter) (*tusd.Handler, error) {

	tusd.NewStoreComposer()
	store := s3store.New(bucketName, adapter.S3())

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:      "/uploads/",
		StoreComposer: composer,
		//PreUploadCreateCallback: func(hook tusd.HookEvent) error {
		//	targetBucket := hook.HTTPRequest.Header.Get("Bucket")
		//	if targetBucket == "" {
		//		return errors.New("bucket name is required")
		//	}
		//	return v.ValidateUpload(targetBucket, hook.Upload.Size)
		//},
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

func NewHandlerWithFileStore(bucketName string, v UploadValidator) (*tusd.Handler, error) {

	store := filestore.New("./uploads/" + bucketName)

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
	if _, err := os.Stat("./uploads/" + bucketName); os.IsNotExist(err) {
		if err := os.Mkdir("./uploads/"+bucketName, os.ModePerm); err != nil {
			log.Fatalf("Could not create upload dir: %v", err)
		}
	}

	if err != nil {
		return nil, err
	}

	return handler, err
}
