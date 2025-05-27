package service

import (
	"github.com/gocastsian/roham/vectorlayerapp/job"
	"go.temporal.io/sdk/temporal"
	"time"

	"go.temporal.io/sdk/workflow"
)

type Workflow struct {
	service Service
}

func New(service Service) Workflow {
	return Workflow{service: service}
}

func (w Workflow) ImportLayerWorkflow(ctx workflow.Context, event job.Event) error {
	fileKey, ok := event.Args["key"].(string)
	if !ok {
		_ = workflow.ExecuteActivity(ctx, w.service.UpdateJob, UpdateJobStatusRequest{
			WorkflowId: event.WorkflowId,
			Status:     JobStatusFailed,
		}).Get(ctx, nil)
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Hour * 24,
		HeartbeatTimeout:       time.Minute * 5,
		ScheduleToCloseTimeout: time.Hour * 24,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	err := workflow.ExecuteActivity(ctx, w.service.UpdateJob, UpdateJobStatusRequest{
		WorkflowId: event.WorkflowId,
		Status:     JobStatusProcessing,
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to update job Status", "Error", err)
		return err
	}

	var importResult ImportLayerResponse
	err = workflow.ExecuteActivity(ctx, w.service.ImportLayer, ImportLayerRequest{
		FileKey: fileKey,
	}).Get(ctx, &importResult)
	if err != nil {
		_ = workflow.ExecuteActivity(ctx, w.service.UpdateJob, UpdateJobStatusRequest{
			WorkflowId: event.WorkflowId,
			Status:     JobStatusFailed,
			ErrorMsg:   err.Error(),
		}).Get(ctx, nil)

		_ = workflow.ExecuteActivity(ctx, w.service.SendNotification, SendNotificationRequest{
			WorkflowId: event.WorkflowId,
			Status:     "failed",
		}).Get(ctx, nil)
		logger.Error("Failed to import layer", "Error", err)
		return err
	}

	var createLayer CreateLayerResponse
	err = workflow.ExecuteActivity(ctx, w.service.CreateLayer, CreateLayerRequest{LayerName: importResult.LayerName}).
		Get(ctx, &createLayer)
	if err != nil {
		var dropTable DropLayerResponse
		_ = workflow.ExecuteActivity(ctx, w.service.DropLayerTable, DropLayerRequest{TableName: importResult.LayerName}).Get(
			ctx, &dropTable)

		_ = workflow.ExecuteActivity(ctx, w.service.UpdateJob, UpdateJobStatusRequest{
			WorkflowId: event.WorkflowId,
			Status:     JobStatusFailed,
			ErrorMsg:   err.Error(),
		})
		_ = workflow.ExecuteActivity(ctx, w.service.SendNotification, SendNotificationRequest{
			WorkflowId: event.WorkflowId,
			Status:     "failed",
		})
		logger.Error("Failed to create layer", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx, w.service.UpdateJob, UpdateJobStatusRequest{
		WorkflowId: event.WorkflowId,
		Status:     JobStatusComplete,
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to update job Status", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx, w.service.SendNotification, UpdateJobStatusRequest{
		WorkflowId: event.WorkflowId,
		Status:     JobStatusComplete,
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send notification", "Error", err)
	}

	return nil
}
