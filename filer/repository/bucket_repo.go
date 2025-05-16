package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/types"
	"log/slog"
)

type BucketRepo struct {
	Logger *slog.Logger
	db     *sql.DB
}

func NewBucketRepo(logger *slog.Logger, db *sql.DB) BucketRepo {
	return BucketRepo{
		Logger: logger,
		db:     db,
	}
}

func (r BucketRepo) CreateBucket(ctx context.Context, bucketName string) (types.ID, error) {
	query := `INSERT INTO filer_buckets (bucket_name) VALUES ($1) RETURNING id`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id types.ID
	err = stmt.QueryRowContext(ctx, bucketName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create bucket: %w", err)
	}

	return id, nil
}

func (r BucketRepo) BucketIsExist() {

}
