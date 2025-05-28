package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/vectorlayerapp/service"
	"strings"
)

type JobRepo struct {
	PostgreSQL *sql.DB // PostgreSQL connection
}

func NewJobRepo(db *sql.DB) LayerRepo {
	return LayerRepo{
		PostgreSQL: db,
	}
}

func (r LayerRepo) AddJob(ctx context.Context, job service.JobEntity) (types.ID, error) {
	query := `INSERT INTO jobs(token, status) VALUES ($1 , $2) returning id;`
	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var res int64
	err = stmt.QueryRowContext(ctx, job.Token, job.Status).Scan(&res)
	if err != nil {
		return 0, err
	}

	return types.ID(res), nil
}

func (r LayerRepo) GetJobByToken(ctx context.Context, token string) (service.JobEntity, error) {
	query := `SELECT * FROM jobs WHERE token = $1;`
	var (
		job service.JobEntity
		id  int64
	)

	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return job, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, token).Scan(&id, &job.Token, &job.Status, &job.Error, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return job, err
	}

	job.ID = types.ID(id)
	return job, nil

}
func (r LayerRepo) UpdateJob(ctx context.Context, job service.JobEntity) (bool, error) {
	setParts := []string{}
	args := []interface{}{}
	argIdx := 1

	if job.Status != "" {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, job.Status)
		argIdx++
	}

	if job.Error != nil {
		setParts = append(setParts, fmt.Sprintf("error = $%d", argIdx))
		args = append(args, job.Error)
		argIdx++
	}

	if len(setParts) == 0 {
		return false, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE jobs SET %s WHERE token = $%d", strings.Join(setParts, ", "), argIdx)
	args = append(args, job.Token)

	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return false, err
	}

	return true, nil
}
