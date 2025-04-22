package importlayer

import "github.com/gocastsian/roham/types"

type CreateJobRequest struct {
	Token    string   `json:"token"`
	FileData []byte   `json:"fileData"`
	UserId   types.ID `json:"user_id"`
}

type CreateJobResponse struct {
	WorkflowID string `json:"workflow_id"`
}
