package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/types"

	"log/slog"
)

type FileMetadataRepo struct {
	Logger *slog.Logger
	db     *sql.DB
}

func NewFileMetadataRepo(logger *slog.Logger, db *sql.DB) FileMetadataRepo {
	return FileMetadataRepo{
		Logger: logger,
		db:     db,
	}
}

func (r FileMetadataRepo) Create(ctx context.Context, f storage.FileMetadata) (types.ID, error) {

	query := `INSERT INTO filer_metadata (storage_key, bucket_name, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id types.ID
	err = stmt.QueryRowContext(ctx, f.Key, f.BucketName, f.Metadata, f.CreatedAt, f.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert file record: %w", err)
	}

	return id, nil

}
