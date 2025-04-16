package user_test

import (
	"context"
	"testing"

	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/guard"
	"github.com/gocastsian/roham/userapp/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAllUsers(ctx context.Context) ([]user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user.User, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockRepository) RegisterUser(ctx context.Context, u user.User) (types.ID, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(types.ID), args.Error(1)
}

func (m *MockRepository) CheckUserUniquness(ctx context.Context, email string, username string) (bool, error) {
	args := m.Called(ctx, email, username)
	return args.Bool(0), args.Error(1)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	userValidator := user.NewValidator(mockRepo)

	guardSvc := &guard.Service{}
	service := user.NewService(mockRepo, userValidator, nil, guardSvc)

	regReq := user.RegisterRequest{
		Username:        "testuser",
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@example.com",
		Password:        "s2Securepassword",
		ConfirmPassword: "s2Securepassword",
	}

	mockRepo.On("CheckUserUniquness", mock.Anything, regReq.Email, regReq.Username).Return(false, nil)
	mockRepo.On("RegisterUser", mock.Anything, mock.AnythingOfType("user.User")).Return(types.ID(1), nil)

	resp, err := service.RegisterUser(context.Background(), regReq)
	assert.NoError(t, err)
	assert.Equal(t, types.ID(0), resp.ID)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	userValidator := user.NewValidator(mockRepo)
	service := user.NewService(mockRepo, userValidator, nil, nil)

	regReq := user.RegisterRequest{
		Username:        "testuser",
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@example.com",
		Password:        "se2Scurepassword",
		ConfirmPassword: "se2Scurepassword",
	}

	mockRepo.On("CheckUserUniquness", mock.Anything, regReq.Email, regReq.Username).Return(true, nil)

	_, err := service.RegisterUser(context.Background(), regReq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user already exist")
	mockRepo.AssertExpectations(t)
}
