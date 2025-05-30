package filestorage

import (
	"time"

	"github.com/gocastsian/roham/types"
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
	FileSize  int64      `json:"file_size"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
