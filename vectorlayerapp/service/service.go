package service

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/types"
	"github.com/google/uuid"
	"log"
	"os/exec"
	"time"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
	AddJob(ctx context.Context, job JobEntity) (types.ID, error)
	GetJobByToken(ctx context.Context, token string) (JobEntity, error)
	UpdateJob(ctx context.Context, job JobEntity) (bool, error)
}

type Scheduler interface {
	Add(ctx context.Context, workflowId string, workflowName string, queueName string) (string, error)
}

type Service struct {
	repository Repository
	validator  Validator
	Scheduler  Scheduler
}

func NewService(repo Repository, validator Validator, scheduler Scheduler) Service {
	return Service{
		repository: repo,
		validator:  validator,
		Scheduler:  scheduler,
	}
}

func (s Service) HealthCheckSrv(ctx context.Context) (string, error) {
	check, err := s.repository.HealthCheck(ctx)
	if err != nil {
		return "", err
	}
	return check, nil
}

func (s Service) ScheduleImportLayer(ctx context.Context) (ScheduleImportLayerResponse, error) {
	workflowId := "layer_" + uuid.New().String()

	_, err := s.repository.AddJob(ctx, JobEntity{
		Token:  workflowId,
		Status: JobStatusPending,
	})
	if err != nil {
		return ScheduleImportLayerResponse{}, fmt.Errorf("failed to create job record: %w", err)
	}

	_, err = s.Scheduler.Add(ctx, workflowId, "ImportLayerWorkflow", "import_layer")
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

// TODO: make this better
func (s Service) ImportLayer(ctx context.Context) (ImportLayerResponse, error) {
	connStr := "PG:host=localhost user=nimamleo dbname=vectorlayer_db password=root"
	shapefile := "/home/nimamleo/Downloads/Iran_shipefile/Ostan.shp"

	cmd := exec.CommandContext(ctx, "ogr2ogr",
		"-f", "PostgreSQL",
		connStr,
		shapefile,
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

	time.Sleep(time.Second * 10) //simulate it takes time

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
