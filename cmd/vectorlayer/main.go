package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gocastsian/roham/pkg/postgresqlmigrator"
	"github.com/gocastsian/roham/vectorlayerapp"

	cfgloader "github.com/gocastsian/roham/pkg/cfg_loader"
	"github.com/gocastsian/roham/pkg/logger"
	"github.com/gocastsian/roham/pkg/postgresql"
)

func main() {
	var cfg vectorlayerapp.Config
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	options := cfgloader.Option{
		Prefix:       "VECTORLAYER_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "vectorlayer", "development", "config.yaml"),
		CallbackEnv:  nil,
	}

	if err := cfgloader.Load(options, &cfg); err != nil {
		log.Fatalf("Failed to load vectorlayerapp config: %v", err)
	}

	logger.Init(cfg.Logger)
	vectorLayerLogger := logger.L()

	vectorLayerLogger.Info("vector layer service started...")

	//todo retry to connect in result of connection failure
	//todo add metrics (each connection)
	postgresConn, cnErr := postgresql.Connect(cfg.PostgresDB)

	if cnErr != nil {
		log.Fatal(cnErr)
	} else {
		vectorLayerLogger.Info(fmt.Sprintf("You are connected to %s successfully.", cfg.PostgresDB.DBName))
	}

	if err != nil {
		log.Fatalf("Error in Connecting to vector layer Postgresql: %v", err)
	}

	mgr := postgresqlmigrator.New(cfg.PostgresDB, cfg.PostgresDB.PathOfMigration)
	mgr.Up()

	defer postgresql.Close(postgresConn.DB)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := vectorlayerapp.Setup(ctx, cfg, postgresConn, vectorLayerLogger)
	app.Start()
}
