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

func (r FileMetadataRepo) Insert(ctx context.Context, i storage.CreateFileMetadataInput) (types.ID, error) {

	query := `
		INSERT INTO file_metadata (
			storage_id,
			file_key,
			file_name,
			mime_type,
			file_size
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id types.ID
	err = stmt.QueryRowContext(ctx,
		i.StorageID,
		i.FileKey,
		i.FileName,
		i.MimeType,
		i.Size,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert file metadata : %w", err)
	}

	return id, nil

}
