package storage

import "github.com/gocastsian/roham/types"

type CreateFileMetadataInput struct {
}

type CreateFileMetadataOutput struct {
	ID         types.ID          `json:"id"`
	Key        string            `json:"key"`
	BucketName string            `json:"bucket_name"`
	Metadata   map[string]string `json:"metadata"`
}

type CreateBucketInput struct {
	BucketName string `json:"bucket_name"`
}

type CreateBucketOutput struct {
	ID   types.ID `json:"bucket_id"`
	Name string   `json:"bucket_name"`
}
