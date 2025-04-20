package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/types"
	"log/slog"

	"github.com/gocastsian/roham/userapp/service/user"
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
	query := `SELECT id, username, first_name, last_name, email, phone_number, birth_date, created_at, updated_at, role FROM users;`

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
			&result.Role,
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

func (repo UserRepo) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user.User, error) {

	// Check if the user exists
	exists, err := repo.checkUserExist(ctx, phoneNumber)
	if !exists || err != nil {
		return user.User{}, fmt.Errorf("failed to check user existence: %w", err)
	}

	// Query to fetch the user's details
	query := "SELECT id, phone_number, role, password_hash FROM users WHERE phone_number=$1"
	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var usr user.User
	err = stmt.QueryRowContext(ctx, phoneNumber).Scan(&usr.ID, &usr.PhoneNumber, &usr.Role, &usr.PasswordHash)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return usr, nil
}

func (repo UserRepo) checkUserExist(ctx context.Context, phoneNumber string) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE phone_number = $1
		)
	`

	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRowContext(ctx, phoneNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute prepared statement: %w", err)
	}

	return true, nil
}

func (repo UserRepo) GetUser(ctx context.Context, ID types.ID) (user.User, error) {
	// Check if the user exists
	exists, err := repo.checkUserExistByID(ctx, ID)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exists {
		return user.User{}, fmt.Errorf("the user not found")
	}

	// Query to fetch the user's details
	query := "SELECT id, username, first_name, last_name, email, phone_number, birth_date, created_at, updated_at, role FROM users WHERE ID=$1"
	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var usr user.User
	err = stmt.QueryRowContext(ctx, ID).Scan(
		&usr.ID,
		&usr.Username,
		&usr.FirstName,
		&usr.LastName,
		&usr.Email,
		&usr.PhoneNumber,
		&usr.BirthDate,
		&usr.CreatedAt,
		&usr.UpdatedAt,
		&usr.Role,
	)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return usr, nil
}

func (repo UserRepo) checkUserExistByID(ctx context.Context, ID types.ID) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE ID = $1
		)
	`

	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRowContext(ctx, ID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute prepared statement: %w", err)
	}

	return exists, nil
}

func (repo UserRepo) UpdateAvatar(ctx context.Context, ID types.ID, uploadAddress string) error {
	query := `UPDATE users SET avatar=$1 WHERE ID=$2`
	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, uploadAddress, ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	if rowsAffected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	} else {
		if rowsAffected != 1 {
			return fmt.Errorf("failed to update the avatar")
		}
	}

	return nil
}
