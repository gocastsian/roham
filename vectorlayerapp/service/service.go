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
	localDir := "./shapefile"

	if _, err := os.Stat(localDir); err == nil {
		if err := os.RemoveAll(localDir); err != nil {
			return ImportLayerResponse{}, fmt.Errorf("failed to remove existing dir %s: %w", localDir, err)
		}
	}

	if err := os.MkdirAll(localDir, 0755); err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to create dir %s: %w", localDir, err)
	}

	data, err := s.queryClient.DownloadShapeFile(fileKey)
	if err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to download %s: %w", fileKey, err)
	}

	zipPath := filepath.Join(localDir, filepath.Base(fileKey)+".zip")
	if err := os.WriteFile(zipPath, data, 0644); err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to write zip file %s: %w", zipPath, err)
	}
	log.Printf("Saved zip file %s", zipPath)

	if err := archiver.Unarchive(zipPath, localDir); err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to unzip file %s: %w", zipPath, err)
	}
	log.Printf("Unzipped files to %s", localDir)

	connStr := "PG:host=localhost user=nimamleo dbname=vectorlayer_db password=root"

	cmd := exec.CommandContext(ctx, "ogr2ogr",
		"-f", "PostgreSQL",
		connStr,
		filepath.Join(localDir, "/ostan/Ostan.shp"),
		"-nln", "ostan",
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
