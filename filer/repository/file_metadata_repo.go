package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gocastsian/roham/filer/service/filestorage"
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

func (r FileMetadataRepo) InsertFileMetadata(ctx context.Context, fileMetadata filestorage.FileMetadata) (types.ID, error) {

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
		fileMetadata.StorageID,
		fileMetadata.FileKey,
		fileMetadata.FileName,
		fileMetadata.MimeType,
		fileMetadata.FileSize,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert file metadata : %w", err)
	}

	return id, nil

}

func (r FileMetadataRepo) FindByKey(ctx context.Context, key string) (filestorage.FileMetadata, error) {
	query := `
        SELECT id, storage_id, file_key, file_name, mime_type, file_size,created_at, claimed_at
        FROM file_metadata WHERE file_key = $1
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return filestorage.FileMetadata{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var f filestorage.FileMetadata
	err = stmt.QueryRowContext(ctx, key).Scan(
		&f.ID,
		&f.StorageID,
		&f.FileKey,
		&f.FileName,
		&f.MimeType,
		&f.FileSize,
		&f.CreatedAt,
		&f.ClaimedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return filestorage.FileMetadata{}, errors.New("file metadata not found")
		}
		return filestorage.FileMetadata{}, fmt.Errorf("failed to find file metadata by Key: %w", err)
	}

	return f, nil
}
