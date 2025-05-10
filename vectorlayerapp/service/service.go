package service

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/vectorlayerapp/job"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
	AddJob(ctx context.Context, job JobEntity) (types.ID, error)
	GetJobByToken(ctx context.Context, token string) (JobEntity, error)
	UpdateJob(ctx context.Context, job JobEntity) (bool, error)
}

type Scheduler interface {
	Add(ctx context.Context, event job.Event) (string, error)
}

type QueryClient interface {
	DownloadShapeFile(fileKey string) ([]byte, error)
}

type Service struct {
	repository  Repository
	validator   Validator
	scheduler   Scheduler
	queryClient QueryClient
}

func NewService(repo Repository, validator Validator, scheduler Scheduler, queryClient QueryClient) Service {
	return Service{
		repository:  repo,
		validator:   validator,
		scheduler:   scheduler,
		queryClient: queryClient,
	}
}

func (s Service) HealthCheckSrv(ctx context.Context) (string, error) {
	check, err := s.repository.HealthCheck(ctx)
	if err != nil {
		return "", err
	}
	return check, nil
}

func (s Service) ScheduleImportLayer(ctx context.Context, fileKey string) (ScheduleImportLayerResponse, error) {
	workflowId := "layer_" + uuid.New().String()

	_, err := s.repository.AddJob(ctx, JobEntity{
		Token:  workflowId,
		Status: JobStatusPending,
	})
	if err != nil {
		return ScheduleImportLayerResponse{}, fmt.Errorf("failed to create job record: %w", err)
	}

	_, err = s.scheduler.Add(ctx, job.Event{
		WorkflowId:   workflowId,
		WorkflowName: "ImportLayerWorkflow",
		QueueName:    "import_layer",
		Args: map[string]any{
			"key": fileKey},
	})

	if err != nil {
		_, _ = s.repository.AddJob(ctx, JobEntity{
			Token:  workflowId,
			Status: JobStatusFailed,
		})
		return ScheduleImportLayerResponse{}, fmt.Errorf("failed to start workflow: %w", err)
	}

	return ScheduleImportLayerResponse{
		WorkflowId: workflowId,
	}, nil
}

func (s Service) UpdateJobStatus(ctx context.Context, req UpdateJobStatusRequest) error {
	_, err := s.repository.UpdateJob(ctx, JobEntity{
		Token:  req.WorkflowId,
		Status: req.Status,
	})
	if err != nil {
		return fmt.Errorf("failed to update job Status: %w", err)
	}
	return err
}

func (s Service) ImportLayer(ctx context.Context, fileKey string) (ImportLayerResponse, error) {
	tempDir, err := os.MkdirTemp("", "shapefile-*")
	if err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	log.Printf("Created temporary directory: %s", tempDir)

	data, err := s.queryClient.DownloadShapeFile(fileKey)
	if err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to download %s: %w", fileKey, err)
	}

	zipPath := filepath.Join(tempDir, filepath.Base(fileKey)+".zip")
	if err := os.WriteFile(zipPath, data, 0644); err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to write zip file %s: %w", zipPath, err)
	}
	log.Printf("Saved zip file %s", zipPath)

	if err := archiver.Unarchive(zipPath, tempDir); err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to unzip file %s: %w", zipPath, err)
	}
	log.Printf("Unzipped files to %s", tempDir)

	var shpFilePath string
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".shp" {
			shpFilePath = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to scan directory for .shp files: %w", err)
	}

	if shpFilePath == "" {
		return ImportLayerResponse{}, fmt.Errorf("no .shp file found in the extracted directory")
	}

	log.Printf("Found shapefile: %s", shpFilePath)

	layerName := strings.ToLower(filepath.Base(shpFilePath[:len(shpFilePath)-4]))

	connStr := "PG:host=localhost user=nimamleo dbname=vectorlayer_db password=root"

	cmd := exec.CommandContext(ctx, "ogr2ogr",
		"-f", "PostgreSQL",
		connStr,
		shpFilePath,
		"-nln", layerName,
		"-overwrite",
		"-append",
		"-nlt", "MULTIPOLYGON",
		"-t_srs", "EPSG:4326",
		"-lco", "GEOMETRY_NAME=wkb_geometry",
		"-lco", "FID=ogc_fid",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ogr2ogr failed: %v\nOutput: %s", err, string(output))
		return ImportLayerResponse{}, fmt.Errorf("ogr2ogr failed: %w", err)
	}

	time.Sleep(10 * time.Second)

	log.Println("Shapefile imported successfully!")
	return ImportLayerResponse{
		Status: true,
	}, nil
}

// TODO: implement real notification
func (s Service) SendNotification(ctx context.Context, req SendNotificationRequest) error {

	_, err := s.repository.GetJobByToken(ctx, req.WorkflowId)
	if err != nil {
		return fmt.Errorf("failed to get job info: %w", err)
	}

	message := fmt.Sprintf(
		"\n=== IMPORT JOB NOTIFICATION ===\n"+
			"Workflow ID: %s\n"+
			"JobStatus: %s\n"+
			"Completion Time: %s\n"+
			"===============================\n",
		req.WorkflowId,
		req.Status,
		time.Now(),
	)

	log.Println(message)

	return nil
}
