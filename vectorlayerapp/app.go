package vectorlayerapp

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/adapter/temporal"
	job "github.com/gocastsian/roham/vectorlayerapp/job/temporal"
	"github.com/gocastsian/roham/vectorlayerapp/service"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/postgresql"
	"github.com/gocastsian/roham/vectorlayerapp/delivery/http"
	"github.com/gocastsian/roham/vectorlayerapp/repository"
)

type Application struct {
	layerRepo  service.Repository
	layerSrv   service.Service
	Handler    http.Handler
	HTTPServer http.Server
	Temporal   temporal.Adapter
	Config     Config
	Logger     *slog.Logger
}

func Setup(ctx context.Context, config Config, postgresConn *postgresql.Database, logger *slog.Logger) Application {

	LayerRepo := repository.NewLayerRepo(postgresConn.DB)
	LayerValidator := service.NewValidator(LayerRepo)
	LayerSrv := service.NewService(LayerRepo, LayerValidator)
	temporalAdp := temporal.New(config.Temporal)
	Handler := http.NewHandler(LayerSrv, logger, temporalAdp)

	return Application{
		layerRepo:  LayerRepo,
		layerSrv:   LayerSrv,
		Handler:    Handler,
		HTTPServer: http.New(httpserver.New(config.HTTPServer), Handler, logger),
		Config:     config,
		Logger:     logger,
		Temporal:   temporalAdp,
	}
}

func (app Application) Start() {
	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startServers(app, &wg)
	startWorkers(app, &wg)

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

	wg.Wait()
	app.Logger.Info("user_app stopped")
}

func startServers(app Application, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info(fmt.Sprintf("HTTP server started on %d", app.Config.HTTPServer.Port))
		if err := app.HTTPServer.Serve(); err != nil {
			// todo add metrics
			app.Logger.Error(
				fmt.Sprintf("error in HTTP server on %d", app.Config.HTTPServer.Port),
				slog.Any("err", err),
			)
		}
		app.Logger.Info(fmt.Sprintf("HTTP server stopped %d", app.Config.HTTPServer.Port))
	}()
}

func startWorkers(app Application, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		worker := job.New(app.Temporal.Client, "greeting")

		worker.RegisterWorkflow(app.layerSrv.HealthCheckJob)
		worker.RegisterActivity(app.layerRepo.HealthCheckJob)

		if err := worker.Start(); err != nil {
			log.Fatalf("error in running worker with err: %v", err)
		}
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
	httpShutdownCtx, httpCancel := context.WithTimeout(context.Background(), app.Config.HTTPServer.ShutDownCtxTimeout)
	defer httpCancel()
	if err := app.HTTPServer.Stop(httpShutdownCtx); err != nil {
		app.Logger.Error(
			fmt.Sprintf("HTTP server graceful shutdown failed: %v", err),
			slog.Any("err", err),
		)
	}
}
