package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/user"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestCheckUserUniquness(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := slog.Default()
	repo := NewUserRepo(Config{}, db, logger)

	// Set expectations for the SQL query
	mock.
		ExpectPrepare("SELECT\\s+EXISTS \\(SELECT 1 FROM users WHERE username = \\$1\\) AS username_exists,\\s+EXISTS \\(SELECT 1 FROM users WHERE email = \\$2\\) AS email_exists").
		ExpectQuery().
		WithArgs("testuser", "test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"username_exists", "email_exists"}).AddRow(true, false))

	ctx := context.Background()
	exists, err := repo.CheckUserUniquness(ctx, "test@example.com", "testuser")

	assert.NoError(t, err)
	assert.True(t, exists)

	// Check expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestRegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	logger := slog.Default()
	repo := NewUserRepo(Config{}, db, logger)
	mock.ExpectPrepare(`INSERT INTO users \(username,first_name,last_name,email,role,password_hash,is_active\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING id`).
		ExpectQuery().
		WithArgs("test", "firstname", "lastname", "email@gmail.com", 0, "password_hash", true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	ctx := context.Background()
	user := user.User{
		ID:           0,
		Username:     "test",
		FirstName:    "firstname",
		LastName:     "lastname",
		Email:        "email@gmail.com",
		Avatar:       "",
		IsActive:     true,
		Role:         0,
		PasswordHash: "password_hash",
	}
	id, err := repo.RegisterUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, types.ID(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())

}
