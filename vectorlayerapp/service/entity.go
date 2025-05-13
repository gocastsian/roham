package service

import (
	"github.com/gocastsian/roham/types"
	"time"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusComplete   JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type JobEntity struct {
	ID        types.ID  `json:"id"`
	Token     string    `json:"token"`
	Status    JobStatus `json:"Status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TODO add new column system_cordiante
type LayerEntity struct {
	ID           types.ID  `json:"id"`
	Name         string    `json:"name"`
	DefaultStyle string    `json:"default_style"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
