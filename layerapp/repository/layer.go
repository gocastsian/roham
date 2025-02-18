package repository

import (
	"context"
	"database/sql"
	"fmt"
	"roham/layerapp/service/layer"
)

// LayerRepo is the concrete implementation of the service.Repository interface
type LayerRepo struct {
	PostgreSQL *sql.DB // PostgreSQL connection
	//Redis      *redis.Client // Redis client connection
}

// NewLayerRepo creates a new instance of LayerRepo with PostgreSQL and Redis connections
func NewLayerRepo(db *sql.DB /*  ,redis *redis.Client */) LayerRepo {
	return LayerRepo{
		PostgreSQL: db,
		//Redis:      redis,
	}
}

func (r LayerRepo) HealthCheck(ctx context.Context) (string, error) {
	query := `SELECT 1;`

	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf(layer.HealthCheckError.Error(), err)
	}
	defer stmt.Close()

	return "everything is ok", nil
}
