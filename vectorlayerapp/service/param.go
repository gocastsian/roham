package service

import "github.com/gocastsian/roham/types"

type ScheduleImportLayerRequest struct{}
type ScheduleImportLayerResponse struct {
	WorkflowId string
}

// ==========================================================
type UpdateJobStatusRequest struct {
	WorkflowId string
	Status     JobStatus
	ErrorMsg   *string
}
type UpdateJobStatusResponse struct{}

// ==========================================================
type ImportLayerRequest struct {
	FileKey string
}
type ImportLayerResponse struct {
	Status    bool
	LayerName string
}

// ==========================================================

type SendNotificationRequest struct {
	WorkflowId string
	Status     string
}
type SendNotificationResponse struct{}

// ==========================================================
type CreateLayerRequest struct {
	LayerName string
}
type CreateLayerResponse struct {
	ID types.ID
}

// ==========================================================
type DropLayerRequest struct {
	TableName string
}
type DropLayerResponse struct {
	Success bool
}
