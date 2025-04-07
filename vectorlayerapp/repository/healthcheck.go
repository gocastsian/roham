package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/vectorlayerapp/service"
)

// LayerRepo is the concrete implementation of the service.Repository interface
type LayerRepo struct {
	PostgreSQL *sql.DB // PostgreSQL connection
}

// NewLayerRepo creates a new instance of LayerRepo with PostgreSQL and Redis connections
func NewLayerRepo(db *sql.DB) LayerRepo {
	return LayerRepo{
		PostgreSQL: db,
	}
}

func (r LayerRepo) HealthCheck(ctx context.Context) (string, error) {
	query := `SELECT 1;`

	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf(service.HealthCheckError.Error(), err)
	}
	defer stmt.Close()

	return "everything is ok", nil
}

func (r LayerRepo) HealthCheckJob(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("hi %v temporal is ok", name), nil
}
