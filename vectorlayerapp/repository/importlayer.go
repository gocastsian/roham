package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/vectorlayerapp/service/importlayer"
)

type ImportRepo struct {
	PostgreSQL *sql.DB
}

func NewImportRepo(db *sql.DB) ImportRepo {
	return ImportRepo{
		PostgreSQL: db,
	}
}

func (r ImportRepo) HealthCheck(ctx context.Context) (string, error) {
	query := `SELECT 1;`

	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf(importlayer.HealthCheckError.Error(), err)
	}
	defer stmt.Close()

	return "everything is ok hiiiii", nil
}

func (repo ImportRepo) CreateJob(ctx context.Context, job importlayer.Job) (types.ID, error) {
	query := `
		INSERT INTO job (file_key, user_id, workflow_id)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var jobID types.ID
	err = stmt.QueryRowContext(ctx,
		job.Token,
		job.UserID,
		job.WorkflowID,
	).Scan(&jobID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert job: %w", err)
	}

	return jobID, nil
}
