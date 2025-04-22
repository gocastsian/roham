package temporalscheduler

import (
	"context"
	"github.com/gocastsian/roham/adapter/temporal"
	"github.com/gocastsian/roham/vectorlayerapp/service/importlayer"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

type Scheduler struct {
	temporal temporal.Adapter
	service  importlayer.Service
}

func New(temporal temporal.Adapter) Scheduler {
	return Scheduler{
		temporal: temporal,
	}
}

func (w Scheduler) ExecuteImportLayer(ctx context.Context, name string) (string, error) {
	workflowId := "test" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: GREETING_QUEUE_NAME,
	}

	we, err := w.temporal.GetClient().ExecuteWorkflow(ctx, options, w.service.ImportLayerWorkflow, name)
	if err != nil {
		return "", err
	}

	var res string
	if err := we.Get(ctx, &res); err != nil {
		return "", err
	}

	return res, nil
}
