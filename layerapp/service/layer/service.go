package layer

import (
	"context"
)

type Repository interface {
	HealthCheck(ctx context.Context) (string, error)
}

type Service struct {
	repository Repository
	validator  Validator
}

func New(repo Repository, validator Validator) Service {
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
