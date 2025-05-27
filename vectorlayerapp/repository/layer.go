package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/types"
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

func (r LayerRepo) CreateLayer(ctx context.Context, layer service.LayerEntity) (types.ID, error) {
	query := `insert into layers(name , default_style) values($1 , $2) returning id;`

	var id types.ID
	err := r.PostgreSQL.QueryRowContext(ctx, query, layer.Name, layer.DefaultStyle).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create layer: %w", err)
	}

	return id, nil
}

func (r LayerRepo) DropTable(ctx context.Context, tableName string) (bool, error) {
	query := fmt.Sprintf(`drop table if exists %s;`, tableName)
	err := r.PostgreSQL.QueryRowContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to drop table: %w", err)
	}
	return true, nil
}
