package repository

import (
	"context"
	"database/sql"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/vectorlayerapp/service"
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

	err = stmt.QueryRowContext(ctx, token).Scan(&id, &job.Token, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return job, err
	}

	job.ID = types.ID(id)
	return job, nil

}

func (r LayerRepo) UpdateJob(ctx context.Context, job service.JobEntity) (bool, error) {
	query := `UPDATE jobs SET status = $1 , error = $2 WHERE token = $3;`
	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, job.Status, job.Error, job.Token)
	if err != nil {
		return true, err
	}
	return true, nil
}
