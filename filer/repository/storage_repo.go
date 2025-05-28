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

func (r StorageRepo) Insert(ctx context.Context, i filestorage.CreateStorageInput) (types.ID, error) {
	query := `
		INSERT INTO storage (kind, name)
		VALUES ($1, $2)
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

func (r StorageRepo) FindByID(ctx context.Context, id types.ID) (filestorage.Storage, error) {
	query := `
        SELECT id, kind, name
        FROM storages
        WHERE id = $1
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return filestorage.Storage{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var s filestorage.Storage
	err = stmt.QueryRowContext(ctx, id).Scan(
		&s.ID,
		&s.Kind,
		&s.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return filestorage.Storage{}, errors.New("storage not found")
		}
		return filestorage.Storage{}, fmt.Errorf("failed to find storage by ID: %w", err)
	}

	return s, nil
}
