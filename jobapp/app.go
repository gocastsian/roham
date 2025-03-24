package jobapp

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/adapter/temporal"
	"github.com/gocastsian/roham/jobapp/delivery/http"
	"github.com/gocastsian/roham/jobapp/service/job"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Application struct {
	HTTPServer http.Server
	Temporal   temporal.Adapter
	Logger     *slog.Logger
	Config     Config
}

func Setup(config Config, logger *slog.Logger) Application {
	testHandler := http.NewHandler()

	return Application{
		HTTPServer: http.New(httpserver.New(config.HTTPServer), testHandler),
		Temporal:   temporal.New(),
		Logger:     logger,
		Config:     config,
	}
}

func (app Application) Start() {
	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startServers(app, &wg)

	<-ctx.Done()
	app.Logger.Info("Shutdown signal received...")

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), app.Config.TotalShutdownTimeout*time.Second)
	defer cancel()

	if app.shutdownServers(shutdownTimeoutCtx) {
		app.Logger.Info("Servers shut down gracefully")
	} else {
		app.Logger.Warn("Shutdown timed out, exiting application")
		os.Exit(1)
	}
	wg.Wait()
	app.Logger.Info("Application stopped")
}

func startServers(app Application, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Logger.Info(fmt.Sprintf("HTTP server strated on %d", app.Config.HTTPServer.Port))
		if err := app.HTTPServer.Server(); err != nil {
			app.Logger.Error(fmt.Sprintf("error in HTTP server on %d", app.Config.HTTPServer.Port))
		}
		app.Logger.Info(fmt.Sprintf("HTTP server stopped on %d", app.Config.HTTPServer.Port))
	}()

	wg.Add(1)
	go func() {
		worker := job.New(app.Temporal.Client, app.Config.Temporal.GreetingQueueName)

		worker.RegisterWorkflow(job.Greeting)
		worker.RegisterActivity(job.SayHelloInPersian)

		if err := worker.Start(); err != nil {
			app.Logger.Error(fmt.Sprintf("error in running worker with err: %v", err))
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

func (app Application) shutdownHTTPServer(we *sync.WaitGroup) {
	defer we.Done()
	httpShutdownCtx, httpCancel := context.WithTimeout(context.Background(), app.Config.HTTPServer.ShutDownCtxTimeout)
	defer httpCancel()

	if err := app.HTTPServer.Stop(httpShutdownCtx); err != nil {
		app.Logger.Error(fmt.Sprintf("HTTP server graceful shutdown failed :%v", err))
	}
}
