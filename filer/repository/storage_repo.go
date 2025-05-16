package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/filer/service/storage"
	"github.com/gocastsian/roham/types"
	"log/slog"
)

type StorageRepo struct {
	Logger *slog.Logger
	db     *sql.DB
}

func NewStorageRepo(logger *slog.Logger, db *sql.DB) StorageRepo {
	return StorageRepo{
		Logger: logger,
		db:     db,
	}
}

func (r StorageRepo) Insert(ctx context.Context, i storage.CreateStorageInput) (types.ID, error) {
	query := `
		INSERT INTO storage (kind, name)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id types.ID
	err = stmt.QueryRowContext(ctx, i.Kind, i.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create i: %w", err)
	}

	return id, nil
}

func (r StorageRepo) StorageIsExist() {

}
