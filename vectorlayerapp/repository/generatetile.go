package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/vectorlayerapp/service/generatetile"
)

type TileRepo struct {
	PostgreSQL *sql.DB
}

func NewTileRepo(db *sql.DB) TileRepo {
	return TileRepo{
		PostgreSQL: db,
	}
}

func (r TileRepo) HealthCheck(ctx context.Context) (string, error) {
	query := `SELECT 1;`

	stmt, err := r.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf(generatetile.HealthCheckError.Error(), err)
	}
	defer stmt.Close()

	return "everything is ok", nil
}
