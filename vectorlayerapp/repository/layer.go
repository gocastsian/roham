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
	query := `insert into layers(name , default_style ,geom_type ) values($1 , $2 , $3) returning id;`

	var id types.ID
	err := r.PostgreSQL.QueryRowContext(ctx, query, layer.Name, layer.DefaultStyle, layer.GeomType).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create layer: %w", err)
	}

	return id, nil
}

func (r LayerRepo) DropTable(ctx context.Context, tableName string) (bool, error) {
	query := `drop table if exists $1;`
	err := r.PostgreSQL.QueryRowContext(ctx, query, tableName)
	if err != nil {
		return false, fmt.Errorf("failed to drop table: %w", err)
	}
	return true, nil
}

func (r LayerRepo) GetLayerByName(ctx context.Context, name string) (service.LayerEntity, error) {
	query := `select * from layers where name = $1;`

	var layer service.LayerEntity
	err := r.PostgreSQL.QueryRowContext(ctx, query, name).Scan(&layer.ID, &layer.Name, &layer.DefaultStyle, &layer.CreatedAt, &layer.UpdatedAt)
	if err != nil {
		return service.LayerEntity{}, fmt.Errorf("failed to read layer %s: %w", name, err)
	}
	return layer, nil
}

func (r LayerRepo) CreateStyle(ctx context.Context, style service.StyleEntity) (types.ID, error) {
	query := `insert into styles(file_path) values($1) returning id;`
	var id types.ID
	err := r.PostgreSQL.QueryRowContext(ctx, query, style.FilePath).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create style %s: %w", style.FilePath, err)
	}
	return id, nil
}
