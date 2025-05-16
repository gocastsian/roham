package command

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gocastsian/roham/filer/adapter/s3adapter"
	"github.com/gocastsian/roham/filer/adapter/tusdadapter"
	"github.com/gocastsian/roham/filer/delivery/http"
	"github.com/gocastsian/roham/filer/delivery/tus"
	"github.com/gocastsian/roham/filer/repository"
	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/filer/service/upload"
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

	awsConfig := aws.NewConfig().
		WithRegion("ir").
		WithEndpoint(cfg.MinioStorage.Endpoint).
		WithCredentials(credentials.NewStaticCredentials(
			cfg.MinioStorage.AccessKey,
			cfg.MinioStorage.SecretKey,
			"", // Leave empty unless using STS/OpenID
		)).
		WithS3ForcePathStyle(true).
		WithDisableSSL(true)

	s3Adapter, err := s3adapter.New(awsConfig)
	if err != nil {
		log.Fatalf("Failed to create AWS session: %s", err)
	}

	postgresConn, err := postgresql.Connect(cfg.PostgresDB)
	if err != nil {
		log.Fatalf("Failed to connect PostgresDB: %s", err)
	}

	mgr := postgresqlmigrator.New(cfg.PostgresDB, cfg.PostgresDB.PathOfMigration)
	mgr.Up()

	defer postgresql.Close(postgresConn.DB)

	fileRepo := repository.NewFileMetadataRepo(appLogger, postgresConn.DB)
	bucketRepo := repository.NewBucketRepo(appLogger, postgresConn.DB)

	// storageService is used for downloading files
	storageService := storage.NewStorageService(s3Adapter, fileRepo, bucketRepo)
	handler := http.NewHandler(storageService)
	httpServer := http.New(httpserver.New(cfg.HTTPServer), handler, appLogger)

	// uploadService will be handling uploading files
	uploadService := upload.NewUploadService(appLogger, fileRepo)

	// todo uploaded files will be moved after upload to default-bucket. should we define different handler per bucket?!
	//tusHandler, err := tusdadapter.NewHandlerWithS3Store("default-bucket", &uploadService, s3Adapter)
	tusHandler, err := tusdadapter.NewHandlerWithFileStore("default-bucket", &uploadService)
	if err != nil {
		log.Fatalf("Failed to create tus handler: %s", err)
	}

	uploadHandler := tus.NewHandler(uploadService, tusHandler)
	uploadServer := tus.NewServer(appLogger, httpserver.New(cfg.Uploader.HTTPServer), uploadHandler)

	app := filer.Setup(ctx, &cfg, appLogger, httpServer, uploadServer, storageService)
	app.Start()
}
