package temporalscheduler

import (
	"context"
	"github.com/gocastsian/roham/adapter/temporal"
	"github.com/gocastsian/roham/vectorlayerapp/job"
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

func (w Scheduler) Add(ctx context.Context, event job.Event) (string, error) {
	options := client.StartWorkflowOptions{
		ID:        event.WorkflowId,
		TaskQueue: event.QueueName,
	}

	we, err := w.temporal.GetClient().ExecuteWorkflow(ctx, options, event.WorkflowName, event)
	if err != nil {
		return "", err
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	return we.GetID(), nil
}
