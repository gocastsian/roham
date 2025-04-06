package main

import (
	"github.com/gocastsian/roham/jobapp"
	cfgloader "github.com/gocastsian/roham/pkg/cfg_loader"
	"github.com/gocastsian/roham/pkg/logger"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var cfg jobapp.Config
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	options := cfgloader.Option{
		Prefix:       "JOB_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workDir, "deploy", "job", "development", "config.yaml"),
		CallbackEnv:  nil,
	}
	if err := cfgloader.Load(options, &cfg); err != nil {
		log.Fatalf("Failed to load jobapp config: %v", err)
	}

	logger.Init(cfg.Logger)
	jobLogger := logger.L()
	jobLogger.Info("job_app service started...")

	app := jobapp.Setup(cfg, jobLogger)
	app.Start()
}
