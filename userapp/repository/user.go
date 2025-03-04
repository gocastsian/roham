package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"roham/userapp/service/user"
)

type Config struct {
}

type UserRepo struct {
	Config     Config
	Logger     *slog.Logger
	PostgreSQL *sql.DB
}

func NewUserRepo(config Config, db *sql.DB, logger *slog.Logger) user.Repository {
	return &UserRepo{
		Config:     config,
		Logger:     logger,
		PostgreSQL: db,
	}
}

func (repo UserRepo) GetAllUsers(ctx context.Context) ([]user.User, error) {
	query := `SELECT id, username, first_name, last_name, email, phone_number, birth_date, created_at, updated_at FROM users;`

	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare find result statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	users := make([]user.User, 0)

	for rows.Next() {
		var result user.User
		var birthDate sql.NullString

		err := rows.Scan(
			&result.ID,
			&result.Username,
			&result.FirstName,
			&result.LastName,
			&result.Email,
			&result.PhoneNumber,
			&birthDate,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}

		if birthDate.Valid {
			result.BirthDate = birthDate.String
		}

		users = append(users, result)
	}

	// Check for any error occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no users found")
	}

	return users, nil
}
