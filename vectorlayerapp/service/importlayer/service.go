package importlayer

import (
	"context"
	"fmt"
	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/statuscode"
	"github.com/gocastsian/roham/types"
	"go.temporal.io/sdk/workflow"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
	CreateJob(ctx context.Context, job Job) (types.ID, error)
}

type Scheduler interface {
	ExecuteImportLayer(ctx context.Context, name string) (string, error)
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

func (s Service) ScheduleImportLayer(ctx context.Context, name string) (string, error) {
	res, err := s.Scheduler.ExecuteImportLayer(ctx, name)

	if err != nil {
		return "", err
	}

	return res, nil
}

func (s Service) ImportLayerWorkflow(ctx workflow.Context, name string) (string, error) {
	return fmt.Sprintf("hi %s", name), nil
}

func (s Service) CreateJob(ctx context.Context, req CreateJobRequest) (CreateJobResponse, error) {
	if err := s.validator.ValidateJobRequest(req); err != nil {
		return CreateJobResponse{}, errmsg.ErrorResponse{
			Message: "Validation error",
			Errors: map[string]interface{}{
				"validation_error": err.Error(),
			},
			InternalErrCode: statuscode.IntCodeInvalidParam,
		}
	}

	workflowID, err := s.Scheduler.ExecuteImportLayer(ctx, "import-layer-workflow-"+req.Token)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("failed to start workflow: %w", err)
	}
	job := Job{
		Token:      req.Token,
		UserID:     req.UserId,
		WorkflowID: workflowID,
	}
	_, err = s.repository.CreateJob(ctx, job)
	if err != nil {
		return CreateJobResponse{}, errmsg.ErrorResponse{
			Message: "Failed to create job",
			Errors: map[string]interface{}{
				"repository_error": err.Error(),
			},
			InternalErrCode: statuscode.IntCodeUnExpected,
		}
	}

	return CreateJobResponse{
		WorkflowID: workflowID,
	}, nil
}
