package jobtemporal

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/adapter/temporal"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type WorkFlow struct {
	temporal temporal.Adapter
}

func New(temporal temporal.Adapter) WorkFlow {
	return WorkFlow{
		temporal: temporal,
	}
}

func (w WorkFlow) HealthCheck(ctx context.Context, name string) (string, error) {
	workflowId := "test" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: GREETING_QUEUE_NAME,
	}

	we, err := w.temporal.GetClient().ExecuteWorkflow(ctx, options, w.HealthCheckWorkflow, name)
	if err != nil {
		return "", err
	}

	var res string
	if err := we.Get(ctx, &res); err != nil {
		return "", err
	}

	return res, nil
}

func (w WorkFlow) HealthCheckWorkflow(ctx workflow.Context, name string) (string, error) {
	return fmt.Sprintf("hi %s\n", name), nil
}
