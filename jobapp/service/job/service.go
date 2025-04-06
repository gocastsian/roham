package job

import (
	"context"
	"fmt"
	temporaladapter "github.com/gocastsian/roham/adapter/temporal"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type Repository interface {
	SayHelloInPersian(ctx context.Context, name string) (string, error)
}
type Service struct {
	Temporal   temporaladapter.Adapter
	Repository Repository
	Config     temporaladapter.Config
}

func NewSvc(temporalAdapter temporaladapter.Adapter, jobRepo Repository, config temporaladapter.Config) Service {
	return Service{
		Temporal:   temporalAdapter,
		Repository: jobRepo,
		Config:     config,
	}
}

func (s Service) Greeting(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Duration(s.Config.StartToCloseTimeout) * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Duration(s.Config.InitialInterval) * time.Second,
			BackoffCoefficient: float64(s.Config.BackoffCoefficient),
			MaximumInterval:    time.Duration(s.Config.MaximumInterval) * time.Second,
			MaximumAttempts:    int32(s.Config.MaximumAttempts),
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var res string
	err := workflow.ExecuteActivity(ctx, s.Repository.SayHelloInPersian, name).Get(ctx, &res)
	if err != nil {
		return "", fmt.Errorf("Failed to get location: %s", err)
	}

	return res, nil
}
