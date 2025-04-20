package userMock

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/user"
	"time"
)

const RepoErr = "repository error"

type UserRepoMock struct {
	users []user.User
}

func DefaultUser() user.User {
	newUser := user.User{
		ID:           1,
		Username:     "default_user",
		FirstName:    "default_firstName",
		LastName:     "default_lastName",
		PhoneNumber:  "",
		Email:        "default@email.com",
		Avatar:       "",
		BirthDate:    "",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		IsActive:     true,
		Role:         1,
		PasswordHash: "",
	}
	return newUser
}

func NewUserRepoMock() *UserRepoMock {
	var users []user.User
	users = append(users, DefaultUser())
	return &UserRepoMock{users: users}
}

func (m *UserRepoMock) GetAllUsers(ctx context.Context) ([]user.User, error) {
	return m.users, nil
}
func (m *UserRepoMock) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user.User, error) {
	for _, u := range m.users {
		if u.PhoneNumber == phoneNumber {
			return u, nil
		}
	}
	return user.User{}, fmt.Errorf("user not found")
}
func (m *UserRepoMock) RegisterUser(ctx context.Context, user user.User) (types.ID, error) {
	lastId := m.users[len(m.users)-1].ID
	user.ID = lastId + 1
	user.IsActive = true
	m.users = append(m.users, user)
	return user.ID, nil
}
func (m *UserRepoMock) CheckUserUniquness(ctx context.Context, email string, username string) (bool, error) {
	for _, u := range m.users {
		if u.Username == username || u.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *UserRepoMock) GetUser(ctx context.Context, ID types.ID) (user.User, error) {
	for _, u := range m.users {
		if u.ID == ID {
			return u, nil
		}
	}
	return user.User{}, nil
}
