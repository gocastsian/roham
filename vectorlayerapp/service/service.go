package service

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/workflow"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
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
