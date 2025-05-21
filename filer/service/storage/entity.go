package storage

import (
	"github.com/gocastsian/roham/types"
	"time"
)

type Storage struct {
	ID   types.ID `json:"id"`
	Name string   `json:"name"`
	Kind string   `json:"kind"`
}

type FileMetadata struct {
	ID        types.ID   `json:"id"`
	StorageID types.ID   `json:"storage_id"`
	FileKey   string     `json:"file_key"`
	FileName  string     `json:"file_name"`
	MimeType  string     `json:"mime_type"`
	Size      string     `json:"file_size"`
	CreatedAt time.Time  `json:"created_at"`
	ClaimedAt *time.Time `json:"claimed_at"`
}

func (f FileMetadata) IsClaimed() bool {
	return f.ClaimedAt != nil
}
