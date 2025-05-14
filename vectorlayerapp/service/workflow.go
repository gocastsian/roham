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

	err := workflow.ExecuteActivity(ctx, w.service.UpdateJobStatus, UpdateJobStatusRequest{
		WorkflowId: event.WorkflowId,
		Status:     JobStatusProcessing,
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to update job Status", "Error", err)
		return err
	}

	var importResult ImportLayerResponse
	err = workflow.ExecuteActivity(ctx, w.service.ImportLayer, event.Args["key"]).Get(ctx, &importResult)
	if err != nil {
		_ = workflow.ExecuteActivity(ctx, w.service.UpdateJobStatus, UpdateJobStatusRequest{
			WorkflowId: event.WorkflowId,
			Status:     JobStatusFailed,
		}).Get(ctx, nil)
		_ = workflow.ExecuteActivity(ctx, w.service.SendNotification, SendNotificationRequest{
			WorkflowId: event.WorkflowId,
			Status:     "failed",
		}).Get(ctx, nil)
		logger.Error("Failed to import layer", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx, w.service.UpdateJobStatus, UpdateJobStatusRequest{
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
