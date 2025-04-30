package command

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gocastsian/roham/filer"
	"github.com/gocastsian/roham/filer/adapter/s3adapter"
	"github.com/gocastsian/roham/filer/upload"
	cfgloader "github.com/gocastsian/roham/pkg/cfg_loader"
	"github.com/gocastsian/roham/pkg/logger"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var serveUploadAppCmd = &cobra.Command{
	Use:   "upload",
	Short: "Start the upload server",
	Long:  `Start the upload server`,
	Run: func(cmd *cobra.Command, args []string) {
		serveUploadApp()
	},
}

func serveUploadApp() {

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

	appLogger.Info("Starting upload Service...")

	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %s", err)
	}

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

	app := upload.Setup(&cfg, appLogger, s3Adapter)
	app.Start()
}
