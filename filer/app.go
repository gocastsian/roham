package filer

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/filer/delivery/http"
	"github.com/gocastsian/roham/filer/delivery/tus"
	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/types"

	"github.com/gocastsian/roham/pkg/postgresql"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Application struct {
	ShutdownCtx  context.Context
	Config       *Config
	Logger       *slog.Logger
	storageSvc   storage.Service
	httpHandler  http.Handler
	HTTPServer   http.Server
	UploadServer tus.Server
	postgresConn *postgresql.Database
}

func Setup(ctx context.Context, cfg *Config, logger *slog.Logger, httpServer http.Server, uploadServer tus.Server, storageSvc storage.Service) Application {
	return Application{
		ShutdownCtx:  ctx,
		Config:       cfg,
		Logger:       logger,
		storageSvc:   storageSvc,
		HTTPServer:   httpServer,
		UploadServer: uploadServer,
	}
}

func (app Application) Start() {

	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// start http server for download of files
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	app.Logger.Info(fmt.Sprintf("HTTP server started on %d", app.Config.HTTPServer.Port))
	//	if err := app.HTTPServer.Serve(); err != nil {
	//		app.Logger.Error(fmt.Sprintf("error in HTTP server on %d", app.Config.HTTPServer.Port), err)
	//	}
	//	app.Logger.Info(fmt.Sprintf("HTTP server stopped %d", app.Config.HTTPServer.Port))
	//}()

	// start http server for download of files
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info(fmt.Sprintf("Upload server started on %d", app.Config.Uploader.HTTPServer.Port))
		if err := app.UploadServer.Serve(); err != nil {
			app.Logger.Error(fmt.Sprintf("error in HTTP server on %d", app.Config.Uploader.HTTPServer.Port), err)
		}
		app.Logger.Info(fmt.Sprintf("HTTP server stopped %d", app.Config.Uploader.HTTPServer.Port))
	}()

	// handle events after uploads
	wg.Add(1)
	go func() {
		for {
			select {
			case info := <-app.UploadServer.EventHandler().CreatedUploads:
				// todo we use some metrics after CreatedUploads
				fmt.Printf("Upload created: %+v\n", info)
			case info := <-app.UploadServer.EventHandler().CompleteUploads:

				input := storage.CreateFileMetadataInput{
					StorageID: 1,
					FileKey:   info.Upload.ID,
					FileName:  info.Upload.MetaData["filename"],
					MimeType:  info.Upload.MetaData["filetype"],
					Size:      info.Upload.Size,
				}
				go func(ctx context.Context, input storage.CreateFileMetadataInput) {
					err := app.UploadServer.Handler.UploadService.OnCompletedUploads(ctx, input)
					if err != nil {
						app.Logger.Error("Unable to handle OnCompletedUploads: %s", err.Error())

					}
				}(ctx, input)
			}
		}
	}()

	go func() {
		for event := range app.UploadServer.Handler.TusHandler.CompleteUploads {
			var bucketName string
			if bucket, ok := event.Upload.MetaData["X-TARGET-STORAGE"]; ok {
				bucketName = bucket
			} else {
				app.Logger.Error("No TARGET-STORAGE specified in header")
				continue
			}

			go func(ctx context.Context, uploadID string, targetStorageName string, metaData map[string]string) {

				input := storage.CreateFileMetadataInput{
					TargetStorageName: targetStorageName,
					FileKey:           uploadID,
					FileName:          metaData["filename"],
					MimeType:          metaData["filetype"],
				}
				err := app.UploadServer.Handler.UploadService.OnCompletedUploads(ctx, input)
				if err != nil {
					app.Logger.Error("Unable to handle OnCompletedUploads: %s", err.Error())
				}
			}(ctx, event.Upload.ID, bucketName, event.Upload.MetaData)
		}
	}()

	<-ctx.Done()
	app.Logger.Info("Shutdown signal received...")

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), app.Config.TotalShutdownTimeout)
	defer cancel()

	if app.shutdownServers(shutdownTimeoutCtx) {
		app.Logger.Info("Servers shut down gracefully")
	} else {
		app.Logger.Warn("Shutdown timed out, exiting application")
		os.Exit(1)
	}

	app.Logger.Info("filer stopped")
}

func (app Application) shutdownServers(ctx context.Context) bool {
	shutdownDone := make(chan struct{})

	go func() {
		var shutdownWg sync.WaitGroup
		shutdownWg.Add(1)
		go app.shutdownHTTPServer(&shutdownWg)

		shutdownWg.Wait()
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		return true
	case <-ctx.Done():
		return false
	}
}

func (app Application) shutdownHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()
	httpShutdownCtx, httpCancel := context.WithTimeout(context.Background(), app.Config.HTTPServer.ShutDownCtxTimeout)
	defer httpCancel()
	if err := app.HTTPServer.Stop(httpShutdownCtx); err != nil {
		app.Logger.Error(fmt.Sprintf("HTTP server graceful shutdown failed: %v", err))
	}
	//todo shutdown upload server
}
