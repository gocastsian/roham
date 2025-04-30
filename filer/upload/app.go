package upload

import (
	"github.com/gocastsian/roham/filer"
	"github.com/gocastsian/roham/filer/adapter/s3adapter"
	"github.com/gocastsian/roham/filer/adapter/tusdadapter"
	tusd "github.com/tus/tusd/pkg/handler"
	tusdS3Store "github.com/tus/tusd/pkg/s3store"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

type Application struct {
	Config        *filer.Config
	Logger        *slog.Logger
	uploadService Service
	s3Adapter     *s3adapter.Adapter
}

func Setup(cfg *filer.Config, logger *slog.Logger, s3adapter *s3adapter.Adapter) Application {
	uploadSvc := NewUploadService(logger)

	return Application{
		Config:        cfg,
		Logger:        logger,
		uploadService: uploadSvc,
		s3Adapter:     s3adapter,
	}
}

func (app Application) Start() {

	app.setupHandler()

	// Start the server
	portStr := ":" + strconv.Itoa(app.Config.Uploader.HTTPServer.Port)
	app.Logger.Info("Starting Uploader Server on port %s ", portStr)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		app.Logger.Error("Unable to start Uploader Server: %s", err)
	}

	//todo handle graceful shutdown
	//todo use echo and middlewares
}

func (app Application) setupHandler() {

	preUploadCreateCallback := func(hook tusd.HookEvent) error {
		req := NewUploadRequest{Size: hook.Upload.Size}
		err := app.uploadService.ValidateUpload(req)
		if err != nil {
			return err
		}
		return nil
	}

	s3 := app.s3Adapter.S3()
	//todo should we setup different store for new buckets or its better to change in on the fly.
	store := tusdS3Store.New("default-bucket", s3)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusdadapter.NewWithS3Store(store, preUploadCreateCallback)

	// we can also use filesystem as tusdStorage.
	//handler, err := tusdadapter.NewWithFileStore(filestore.FileStore{
	//	Path: "./filer/upload/temp",
	//}, preUploadCreateCallback)

	if err != nil {
		log.Fatalf("Failed to create tusd handler: %s", err)
	}

	go func() {
		for event := range handler.CompleteUploads {
			go func() {
				err := app.uploadService.CompleteUpload(event.Upload.ID, event.Upload.MetaData)
				if err != nil {
					app.Logger.Error("Unable to handle CompleteUpload %s", err.Error())
				}
			}()
		}
	}()

	http.Handle("/uploads/", http.StripPrefix("/uploads/", handler))
	http.HandleFunc("/health", healthCheckHandler)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("filer: uploader is running"))
}
