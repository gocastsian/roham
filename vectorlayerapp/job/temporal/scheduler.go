package temporalscheduler

import (
	"context"
	"github.com/gocastsian/roham/adapter/temporal"
	"go.temporal.io/sdk/client"
	"log"
)

type Scheduler struct {
	temporal temporal.Adapter
}

func New(temporal temporal.Adapter) Scheduler {
	return Scheduler{
		temporal: temporal,
	}
}

func (w Scheduler) Add(ctx context.Context, workflowId string, workflowName string, queueName string) (string, error) {
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: queueName,
	}

	we, err := w.temporal.GetClient().ExecuteWorkflow(ctx, options, workflowName, workflowId)
	if err != nil {
		return "", err
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	return we.GetID(), nil
}
