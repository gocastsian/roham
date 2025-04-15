package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/user"
	"log/slog"
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
func (repo UserRepo) CheckUsernameExistence(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)`

	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	var exists bool
	err = stmt.QueryRowContext(ctx, username).Scan(&exists)
	print(exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute prepared statement: %w", err)
	}

	return exists, nil

}
func (repo UserRepo) CheckEmailExistence(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)`
	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	var exists bool
	err = stmt.QueryRowContext(ctx, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute prepared statement: %w", err)
	}

	return exists, nil
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
func (repo UserRepo) CheckUserUniquness(ctx context.Context, email string, username string) (bool, error) {
	query := `SELECT 
				EXISTS (SELECT 1 FROM users WHERE username = $1) AS username_exists,
				EXISTS (SELECT 1 FROM users WHERE email = $2) AS email_exists,
				EXISTS (SELECT 1 FROM users WHERE phone_number = $3) AS phone_number_exists`
	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	var username_exist, email_exist, phonenumber_exist bool
	err = stmt.QueryRowContext(ctx, username, email).Scan(&username_exist, &email_exist, &phonenumber_exist)
	if err != nil {
		return false, fmt.Errorf("failed to execute prepared statement: %w", err)
	}
	print(username_exist, email_exist, phonenumber_exist)
	return username_exist || email_exist || phonenumber_exist, nil
}
func (repo UserRepo) RegisterUser(ctx context.Context, user user.User) (types.ID, error) {
	query := `INSERT INTO users 
    			(username,first_name,last_name,email,role,password_hash) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	stmt, err := repo.PostgreSQL.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id types.ID

	err = stmt.QueryRowContext(ctx,
		user.Username,     // $1: Username
		user.FirstName,    // $2: First name
		user.LastName,     // $3: Last name
		user.Email,        // $4: Email
		user.Role,         // $7: Role
		user.PasswordHash, // $8: Password hash
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to register user: %w", err)
	}

	return id, nil
}
