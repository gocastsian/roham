package filer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gocastsian/roham/filer/adapter/tusdadapter"
	"github.com/gocastsian/roham/filer/delivery/http"
	"github.com/gocastsian/roham/filer/delivery/tus"
	"github.com/gocastsian/roham/filer/service/filestorage"
)

type Application struct {
	ShutdownCtx context.Context
	Config      *Config
	Logger      *slog.Logger
	storageSvc  filestorage.Service
	// httpHandler  http.Handler
	HTTPServer   http.Server
	UploadServer tus.Server
}

func Setup(ctx context.Context, cfg *Config, logger *slog.Logger, httpServer http.Server, uploadServer tus.Server, storageSvc filestorage.Service) Application {
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info(fmt.Sprintf("HTTP server started on %d", app.Config.HTTPServer.Port))
		if err := app.HTTPServer.Serve(); err != nil {
			app.Logger.Info(fmt.Sprintf("error in HTTP server on %d", app.Config.HTTPServer.Port))
		}
		app.Logger.Info(fmt.Sprintf("HTTP server stopped %d", app.Config.HTTPServer.Port))
	}()

	// start http server for download of files
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info(fmt.Sprintf("Upload server started on %d", app.Config.Uploader.HTTPServer.Port))
		if err := app.UploadServer.Serve(); err != nil {
			app.Logger.Error(fmt.Sprintf("error in HTTP server on %d", app.Config.Uploader.HTTPServer.Port))
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

				event := tusdadapter.CompleteUploadsHookEvent(info)
				input, err := event.ConvertToCreateFileMetadataInput()

				if err != nil {
					app.Logger.Error(err.Error())
					return
				}

				go func(ctx context.Context, i filestorage.CreateFileMetadataInput) {
					err := app.UploadServer.Handler.UploadService.OnCompletedUploads(ctx, i)
					if err != nil {
						app.Logger.Error(fmt.Sprintf("Unable to handle OnCompletedUploads: %s", err.Error()))
					}
				}(ctx, input)
			}
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
