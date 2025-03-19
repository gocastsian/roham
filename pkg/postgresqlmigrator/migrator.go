package postgresqlmigrator

import (
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/pkg/postgresql"
	"log"

	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   postgresql.Config
	migrations *migrate.FileMigrationSource
}

func New(dbConfig postgresql.Config, path string) Migrator {

	migrations := &migrate.FileMigrationSource{
		Dir: path,
	}
	return Migrator{dbConfig: dbConfig, dialect: "postgres", migrations: migrations}
}

func (m Migrator) Up() {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		m.dbConfig.User, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DBName, m.dbConfig.SSLMode)

	db, err := sql.Open(m.dialect, connStr)
	if err != nil {
		log.Fatalf("can't open postgres db: %v", err)
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		log.Fatalf("can't apply migrations: %v", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		m.dbConfig.User, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DBName)

	db, err := sql.Open(m.dialect, connStr)
	if err != nil {
		log.Fatalf("can't open postgres db: %v", err)
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		log.Fatalf("can't apply migrations: %v", err)
	}

	fmt.Printf("Rolled back %d migrations!\n", n)
}
