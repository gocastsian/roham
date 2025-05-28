package filestorage

import "github.com/gocastsian/roham/types"

type CreateFileMetadataInput struct {
	TargetStorageID types.ID `json:"target_storage_id"`
	FileKey         string   `json:"file_key"`
	FileName        string   `json:"file_name"`
	MimeType        string   `json:"mime_type"`
	Size            int64    `json:"file_size"`
}

type CreateFileMetadataOutput struct {
	ID       types.ID `json:"id"`
	FileKey  string   `json:"file_key"`
	FileName string   `json:"file_name"`
	MimeType string   `json:"mime_type"`
	Size     int64    `json:"file_size"`
}

type CreateStorageInput struct {
	Name string `json:"storage_name"`
	Kind string `json:"storage_kind"`
}

type CreateStorageOutput struct {
	ID   types.ID `json:"storage_id"`
	Name string   `json:"storage_name"`
	Kind string   `json:"storage_kind"`
}
