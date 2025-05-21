package command

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/filer/adapter/tusdadapter"
	"github.com/gocastsian/roham/filer/delivery/http"
	"github.com/gocastsian/roham/filer/delivery/tus"
	"github.com/gocastsian/roham/filer/repository"
	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/filer/service/upload"
	"github.com/gocastsian/roham/filer/storageprovider/storagefactory"
	cfgloader "github.com/gocastsian/roham/pkg/cfg_loader"
	httpserver "github.com/gocastsian/roham/pkg/http_server"
	"github.com/gocastsian/roham/pkg/logger"
	"github.com/gocastsian/roham/pkg/postgresql"
	"github.com/gocastsian/roham/pkg/postgresqlmigrator"
	"log"
	"os"
	"path/filepath"

	"github.com/gocastsian/roham/filer"
	"github.com/spf13/cobra"
)

var serveFilerCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the filer",
	Long:  `Start the filer`,
	Run: func(cmd *cobra.Command, args []string) {
		serveFiler()
	},
}

func serveFiler() {

	var cfg filer.Config
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v", err)
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "local"
	}

	options := cfgloader.Option{
		Prefix:       "FILER_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "filer", environment, "config.yaml"),
		CallbackEnv:  nil,
	}

	if err := cfgloader.Load(options, &cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.Logger)
	appLogger := logger.L()

	appLogger.Info("Starting filer Service...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresConn, err := postgresql.Connect(cfg.PostgresDB)
	if err != nil {
		log.Fatalf("Failed to connect PostgresDB: %s", err)
	}

	mgr := postgresqlmigrator.New(cfg.PostgresDB, cfg.PostgresDB.PathOfMigration)
	mgr.Up()

	defer postgresql.Close(postgresConn.DB)

	fileRepo := repository.NewFileMetadataRepo(appLogger, postgresConn.DB)
	storageRepo := repository.NewStorageRepo(appLogger, postgresConn.DB)

	storageProvider, err := storagefactory.New(cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to create storage provider: %s", err)
	}

	storageService := storage.NewStorageService(storageProvider, fileRepo, storageRepo)
	handler := http.NewHandler(storageService)
	httpServer := http.New(httpserver.New(cfg.HTTPServer), handler, appLogger)

	// Setup UploadServer
	uploadService := upload.NewUploadService(appLogger, fileRepo)
	tusHandler, err := tusdadapter.New(storageProvider, &uploadService)
	if err != nil {
		log.Fatalf("Failed to create tus handler for storage type %s: %v", cfg.Storage.Type, err)
	}

	uploadHandler := tus.NewHandler(uploadService, tusHandler)
	uploadServer := tus.NewServer(appLogger, httpserver.New(cfg.Uploader.HTTPServer), uploadHandler)

	app := filer.Setup(ctx, &cfg, appLogger, httpServer, uploadServer, storageService)
	app.Start()
}
