package postgresql

import "database/sql"

import (
	_ "github.com/lib/pq"
)

type Adapter struct {
	db *sql.DB
}

func New(connStr string) (*Adapter, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Adapter{db: db}, nil
}

// Exec executes a query that modifies data (e.g., INSERT, UPDATE, DELETE).
func (a *Adapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.db.Exec(query, args...)
}

// Query executes a query that returns rows (e.g., SELECT).
func (a *Adapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return a.db.Query(query, args...)
}

// QueryRow executes a query that returns at most one row.
func (a *Adapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return a.db.QueryRow(query, args...)
}

func (a *Adapter) Close() error {
	return a.db.Close()
}
