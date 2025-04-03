package user

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/types"
	"log/slog"

	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/password"
	"github.com/gocastsian/roham/userapp/service/guard"
)

type Repository interface {
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
	GetUser(ctx context.Context, ID types.ID) (User, error)
}

type Service struct {
	repository Repository
	validator  Validator
	logger     *slog.Logger
	guard      *guard.Service
}

func NewService(repo Repository, validator Validator, logger *slog.Logger, guard *guard.Service) Service {
	return Service{
		repository: repo,
		validator:  validator,
		logger:     logger,
		guard:      guard,
	}
}

func (srv Service) GetAllUsers(ctx context.Context) (GetAllUsersResponse, error) {

	users, err := srv.repository.GetAllUsers(ctx)
	if err != nil {
		srv.logger.Error("user_GetAllUsers", slog.Any("err", err))
		return GetAllUsersResponse{}, errmsg.ErrorResponse{
			Message: err.Error(),
			Errors: map[string]interface{}{
				"user_GetAllUsers": err.Error(),
			},
		}
	}
	responseUsers := make([]GetAllUsersItem, 0)
	for _, user := range users {
		responseUsers = append(responseUsers, GetAllUsersItem{
			ID:          user.ID,
			Username:    user.Username,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Avatar:      user.Avatar,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			BirthDate:   user.BirthDate,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Role:        user.Role,
		})
	}
	return GetAllUsersResponse{Users: responseUsers}, nil
}

func (srv Service) Login(ctx context.Context, loginReq LoginRequest) (LoginResponse, error) {
	if err := srv.validator.ValidatePhoneNumber(loginReq.PhoneNumber); err != nil {
		return LoginResponse{}, errmsg.ErrorResponse{
			Message: "invalid phone number",
			Errors:  map[string]interface{}{"validation_error": err.Error()},
		}
	}

	usr, err := srv.repository.GetUserByPhoneNumber(ctx, loginReq.PhoneNumber)
	if err != nil {
		srv.logger.Error("user_Login", slog.Any("err", err))
		return LoginResponse{}, errmsg.ErrorResponse{
			Message: err.Error(),
			Errors:  map[string]interface{}{"user_Login": err.Error()},
		}
	}

	if !password.CheckPasswordHash(loginReq.Password, usr.PasswordHash) {
		return LoginResponse{}, errmsg.ErrorResponse{
			Message: errmsg.ErrWrongCredentials.Error(),
			Errors:  map[string]interface{}{"user_Login": errmsg.ErrWrongCredentials.Error()},
		}
	}

	userClaim := guard.UserClaim{
		ID:   usr.ID,
		Role: usr.Role,
	}

	accessTok, err := srv.guard.CreateAccessToken(userClaim)
	if err != nil {

		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshTok, err := srv.guard.CreateRefreshToken(userClaim)
	if err != nil {

		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return LoginResponse{
		ID:          usr.ID,
		PhoneNumber: usr.PhoneNumber,
		Tokens: Tokens{
			AccessToken:  accessTok,
			RefreshToken: refreshTok,
		},
	}, nil
}

func (srv Service) GetUser(ctx context.Context, userID types.ID) (GetAllUsersItem, error) {

	user, err := srv.repository.GetUser(ctx, userID)
	if err != nil {
		srv.logger.Error("user_GetUser", slog.Any("err", err))
		return GetAllUsersItem{}, errmsg.ErrorResponse{
			Message: err.Error(),
			Errors: map[string]interface{}{
				"user_GetUser": err.Error(),
			},
		}
	}
	responseUser := GetAllUsersItem{
		ID:          user.ID,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Avatar:      user.Avatar,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		BirthDate:   user.BirthDate,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Role:        user.Role,
	}
	return responseUser, nil
}
