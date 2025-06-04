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
	Error     *string   `json:"Error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TODO add new column system_cordiante
type LayerEntity struct {
	ID           types.ID  `json:"id"`
	Name         string    `json:"name"`
	GeomType     string    `json:"geom_type"`
	DefaultStyle types.ID  `json:"default_style"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StyleEntity struct {
	ID        types.ID  `json:"id"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LayerStylesEntity struct {
	ID        types.ID  `json:"id"`
	LayerID   types.ID  `json:"layer_id"`
	StyleID   types.ID  `json:"style_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
