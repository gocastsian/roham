package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func Connect(config Config) (*Database, error) {
	conn, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	))

	conn.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime))
	conn.SetMaxOpenConns(config.MaxOpenConns)
	conn.SetMaxIdleConns(config.MaxIdleConns)

	return &Database{DB: conn}, err
}

func Close(conn *sql.DB) error {
	return conn.Close()
}

func Ping(conn *sql.DB) error {
	return conn.Ping()
}

func ExampleQuery(db *sql.DB) (string, error) {
	var res string
	err := db.QueryRow("SELECT version()").Scan(&res)
	if err != nil {
		return "", fmt.Errorf("error executing query: %v", err)
	}
	return res, nil
}
