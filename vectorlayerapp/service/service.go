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
	CreateLayer(ctx context.Context, layer LayerEntity) (types.ID, error)
	DropTable(ctx context.Context, tableName string) (bool, error)
	GetLayerByName(ctx context.Context, name string) (LayerEntity, error)
	CreateStyle(ctx context.Context, style StyleEntity) (types.ID, error)
}

type Scheduler interface {
	Add(ctx context.Context, event job.Event) (string, error)
}

type FilerClient interface {
	DownloadShapeFile(fileKey string) ([]byte, error)
}

type Service struct {
	repository  Repository
	validator   Validator
	scheduler   Scheduler
	filerClient FilerClient
}

func NewService(repo Repository, validator Validator, scheduler Scheduler, queryClient FilerClient) Service {
	return Service{
		repository:  repo,
		validator:   validator,
		scheduler:   scheduler,
		filerClient: queryClient,
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

func (s Service) UpdateJob(ctx context.Context, req UpdateJobStatusRequest) error {
	_, err := s.repository.UpdateJob(ctx, JobEntity{
		Token:  req.WorkflowId,
		Status: req.Status,
		Error:  req.ErrorMsg,
	})
	if err != nil {
		return fmt.Errorf("failed to update job Status: %w", err)
	}
	return err
}

func (s Service) ImportLayer(ctx context.Context, req ImportLayerRequest) (ImportLayerResponse, error) {
	tempDir, err := os.MkdirTemp("", "shapefile-*")
	if err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	log.Printf("Created temporary directory: %s", tempDir)

	data, err := s.filerClient.DownloadShapeFile(req.FileKey)
	if err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to download %s: %w", req.FileKey, err)
	}

	zipPath := filepath.Join(tempDir, filepath.Base(req.FileKey)+".zip")
	if err := os.WriteFile(zipPath, data, 0644); err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to write zip file %s: %w", zipPath, err)
	}
	log.Printf("Saved zip file %s", zipPath)

	if err := archiver.Unarchive(zipPath, tempDir); err != nil {
		return ImportLayerResponse{}, fmt.Errorf("failed to unzip file %s: %w", zipPath, err)
	}
	log.Printf("Unzipped files to %s", tempDir)

	var shpFilePath string
	var sldFilePath string
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".shp":
				shpFilePath = path
			case ".sld":
				sldFilePath = path
			}
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

	var styleFileId types.ID
	if sldFilePath != "" {
		log.Printf("Found SLD file: %s", sldFilePath)

		styleRes, err := s.CreateStyle(ctx, CreateStyleRequest{FilePath: sldFilePath})
		if err != nil {
			log.Printf("Warning: Failed to process SLD file: %v", err)
		} else {
			styleFileId = styleRes.ID
			log.Printf("Style created successfully with ID: %v", styleRes)
		}
	}

	log.Println("Shapefile imported successfully!")
	return ImportLayerResponse{
		Status:      true,
		LayerName:   layerName,
		StyleFileID: styleFileId,
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

func (s Service) CreateLayer(ctx context.Context, req CreateLayerRequest) (CreateLayerResponse, error) {
	getLayer, err := s.repository.GetLayerByName(ctx, req.LayerName)
	if err != nil {
		createLayer, err := s.repository.CreateLayer(ctx, LayerEntity{
			Name:         req.LayerName,
			GeomType:     req.GeomType,
			DefaultStyle: req.DefaultStyle,
		})
		if err != nil {
			return CreateLayerResponse{}, fmt.Errorf("failed to create createLayer %s: %w", req.LayerName, err)
		}

		return CreateLayerResponse{
			ID: createLayer,
		}, nil
	}

	return CreateLayerResponse{
		ID: getLayer.ID,
	}, nil
}

func (s Service) DropLayerTable(ctx context.Context, req DropLayerRequest) (DropLayerResponse, error) {
	res, err := s.repository.DropTable(ctx, req.TableName)
	if err != nil {
		return DropLayerResponse{}, fmt.Errorf("failed to drop table %s: %w", req.TableName, err)
	}

	return DropLayerResponse{
		Success: res,
	}, nil

}

func (s Service) CreateStyle(ctx context.Context, req CreateStyleRequest) (CreateStyleResponse, error) {
	stylesDir := "./styles"

	if err := os.MkdirAll(stylesDir, 0755); err != nil {
		return CreateStyleResponse{}, fmt.Errorf("failed to create styles directory: %w", err)
	}

	styleFileName := fmt.Sprintf("style_%s.sld", uuid.New())

	destinationPath := filepath.Join(stylesDir, styleFileName)

	sldContent, err := os.ReadFile(req.FilePath)
	if err != nil {
		return CreateStyleResponse{}, fmt.Errorf("failed to read SLD file %s: %w", req.FilePath, err)
	}

	if err := os.WriteFile(destinationPath, sldContent, 0644); err != nil {
		return CreateStyleResponse{}, fmt.Errorf("failed to write SLD file to %s: %w", destinationPath, err)
	}

	log.Printf("SLD file copied to: %s", destinationPath)

	styleID, err := s.repository.CreateStyle(ctx, StyleEntity{
		FilePath: destinationPath,
	})
	if err != nil {
		if removeErr := os.Remove(destinationPath); removeErr != nil {
			log.Printf("Warning: Failed to clean up file %s after database error: %v", destinationPath, removeErr)
		}
		return CreateStyleResponse{}, fmt.Errorf("failed to create style record in database: %w", err)
	}

	return CreateStyleResponse{ID: styleID}, nil
}
