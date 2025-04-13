package service

import (
	"context"
	jobtemporal "github.com/gocastsian/roham/vectorlayerapp/job/temporal"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
}

type Service struct {
	repository Repository
	validator  Validator
	workflow   jobtemporal.WorkFlow
}

func NewService(repo Repository, validator Validator, workflow jobtemporal.WorkFlow) Service {
	return Service{
		repository: repo,
		validator:  validator,
		workflow:   workflow,
	}
}

func (s Service) HealthCheckSrv(ctx context.Context) (string, error) {
	check, err := s.repository.HealthCheck(ctx)
	if err != nil {
		return "", err
	}
	return check, nil
}

func (s Service) HealthCheckJob(ctx context.Context, name string) (string, error) {
	res, err := s.workflow.Greeting(ctx, name)

	if err != nil {
		return "", err
	}

	return res, nil
}
