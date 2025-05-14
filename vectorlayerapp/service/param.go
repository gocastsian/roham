package service

type ScheduleImportLayerRequest struct{}
type ScheduleImportLayerResponse struct {
	WorkflowId string
}

// ==========================================================
type UpdateJobStatusRequest struct {
	WorkflowId string
	Status     JobStatus
}
type UpdateJobStatusResponse struct{}

// ==========================================================
type ImportLayerRequest struct{}
type ImportLayerResponse struct {
	Status bool
}

// ==========================================================

type SendNotificationRequest struct {
	WorkflowId string
	Status     string
}
type SendNotificationResponse struct{}
