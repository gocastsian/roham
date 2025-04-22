package importlayer

import (
	"github.com/gocastsian/roham/types"
	"time"
)

type Job struct {
	ID           types.ID  `json:"id"`
	Token        string    `json:"file_key"`
	UserID       types.ID  `json:"user_id"`
	WorkflowID   string    `json:"workflow_id"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
