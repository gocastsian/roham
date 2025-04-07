package service

import (
	"context"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
	HealthCheckJob(ctx context.Context, name string) (string, error)
}

type Service struct {
	repository Repository
	validator  Validator
}

func NewService(repo Repository, validator Validator) Service {
	return Service{
		repository: repo,
		validator:  validator,
	}
}

func (s Service) HealthCheckSrv(ctx context.Context) (string, error) {
	check, err := s.repository.HealthCheck(ctx)
	if err != nil {
		return "", err
	}
	return check, nil
}

func (s Service) HealthCheckJob(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 1,
			MaximumInterval:    4 * time.Second,
			MaximumAttempts:    10,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var res string
	err := workflow.ExecuteActivity(ctx, s.repository.HealthCheckJob, name).Get(ctx, &res)

	if err != nil {
		return "", err
	}

	return res, nil
}
