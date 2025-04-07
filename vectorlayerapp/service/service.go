package service

import (
	"context"
	"go.temporal.io/sdk/workflow"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
	HealthCheckJob() (string, error)
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

func (s Service) HealthCheckJob(ctx workflow.Context) (string, error) {
	res, err := s.repository.HealthCheckJob()

	if err != nil {
		return "", err
	}

	return res, nil
}
