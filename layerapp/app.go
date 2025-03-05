package layerapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"roham/layerapp/delivery/http"
	"roham/layerapp/repository"
	"roham/layerapp/service/layer"
	httpserver "roham/pkg/http_server"
	"roham/pkg/postgresql"
)

type Application struct {
	ShutdownCtx  context.Context
	LayerRepo    layer.Repository
	LayerSvc     layer.Service
	LayerHandler http.Handler
	//GRPCServer        grpcServer.Server
	HTTPServer  http.Server
	LayerCfg    Config
	LayerLogger *slog.Logger
}

func Setup(ctx context.Context, config Config, conn *postgresql.Database, logger *slog.Logger) Application {
	layerRepo := repository.NewLayerRepo(conn.DB)
	layerValidator := layer.NewValidator(layerRepo)
	layerSvc := layer.New(layerRepo, layerValidator)
	layerHandler := http.NewHandler(layerSvc, logger)

	return Application{
		ShutdownCtx:  ctx,
		LayerSvc:     layerSvc,
		LayerRepo:    layerRepo,
		LayerHandler: layerHandler,
		HTTPServer:   http.New(httpserver.New(config.Server), layerHandler),
		LayerLogger:  logger,
		LayerCfg:     config,
	}
}
func (app Application) Start() {
	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startServers(app, &wg)
	<-ctx.Done()
	app.LayerLogger.Info("Shutdown signal received...")

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), app.LayerCfg.TotalShutdownTimeout)
	defer cancel()

	if app.shutdownServers(shutdownTimeoutCtx) {
		app.LayerLogger.Info("Servers shut down gracefully")
	} else {
		app.LayerLogger.Warn("Shutdown timed out, exiting application")
		os.Exit(1)
	}

	app.LayerLogger.Info("user_app stopped")
}

func startServers(app Application, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.LayerLogger.Info(fmt.Sprintf("HTTP server started on %d", app.LayerCfg.Server.Port))
		if err := app.HTTPServer.Serve(); err != nil {
			// todo add metrics
			app.LayerLogger.Error(fmt.Sprintf("error in HTTP server on %d", app.LayerCfg.Server.Port), err)
		}
		app.LayerLogger.Info(fmt.Sprintf("HTTP server stopped %d", app.LayerCfg.Server.Port))
	}()
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
	httpShutdownCtx, httpCancel := context.WithTimeout(context.Background(), app.LayerCfg.Server.ShutDownCtxTimeout)
	defer httpCancel()
	if err := app.HTTPServer.Stop(httpShutdownCtx); err != nil {
		app.LayerLogger.Error(fmt.Sprintf("HTTP server graceful shutdown failed: %v", err))
	}
}
