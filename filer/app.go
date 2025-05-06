package filer

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/filer/adapter/s3adapter"
	"github.com/gocastsian/roham/filer/delivery/http"
	"github.com/gocastsian/roham/filer/service/file"
	httpserver "github.com/gocastsian/roham/pkg/http_server"

	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Application struct {
	ShutdownCtx context.Context
	Config      *Config
	Logger      *slog.Logger
	s3Adapter   *s3adapter.Adapter
	httpHandler http.Handler
	HTTPServer  http.Server
}

func Setup(ctx context.Context, cfg *Config, logger *slog.Logger, s3Adapter *s3adapter.Adapter) Application {

	filerService := file.NewFileService(s3Adapter)
	handler := http.NewHandler(filerService)

	httpServer := http.New(httpserver.New(cfg.HTTPServer), handler, logger)

	return Application{
		ShutdownCtx: ctx,
		Config:      cfg,
		Logger:      logger,
		s3Adapter:   s3Adapter,
		httpHandler: handler,
		HTTPServer:  httpServer,
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
			app.Logger.Error(fmt.Sprintf("error in HTTP server on %d", app.Config.HTTPServer.Port), err)
		}
		app.Logger.Info(fmt.Sprintf("HTTP server stopped %d", app.Config.HTTPServer.Port))
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
}
