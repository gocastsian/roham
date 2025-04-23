package file

import "time"

type Info struct {
	ID          string            `json:"id"`
	Key         string            `json:"key"`
	BucketName  string            `json:"bucket_name"`
	Name        string            `json:"name"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
