package user

import (
	"context"
	"fmt"
	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/password"
	"github.com/gocastsian/roham/pkg/statuscode"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/userapp/service/guard"
	"log/slog"
)

type Repository interface {
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
	CheckUsernameExistence(ctx context.Context, username string) (bool, error)
	CheckEmailExistence(ctx context.Context, email string) (bool, error)
	RegisterUser(ctx context.Context, user User) (types.ID, error)
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

func (srv Service) RegisterUser(ctx context.Context, regReq RegisterRequest) (RegisterResponse, error) {
	// validate user registration request fields
	if err := srv.validator.ValidateRegistration(regReq); err != nil {
		return RegisterResponse{}, err
	}
	// check uniqueness of username and email
	if usernameExist, err := srv.repository.CheckUsernameExistence(ctx, regReq.Username); err != nil {
		return RegisterResponse{}, errmsg.ErrorResponse{
			Message: "Application can not detect username existence!",
			Errors:  map[string]interface{}{"user_Register": err.Error()},
		}
	} else if usernameExist {
		return RegisterResponse{}, errmsg.ErrorResponse{
			Message:         "user name exist!",
			InternalErrCode: statuscode.IntCodeUserExistence,
		}
	}

	if emailExist, err := srv.repository.CheckEmailExistence(ctx, regReq.Email); err != nil {
		return RegisterResponse{}, errmsg.ErrorResponse{
			Message: "Application can not detect email existence!",
			Errors:  map[string]interface{}{"user_Register": err.Error()},
		}
	} else if emailExist {
		return RegisterResponse{}, errmsg.ErrorResponse{
			Message:         "email already exist!",
			InternalErrCode: statuscode.IntCodeUserExistence,
		}
	}
	// prepare user entity for save in storage
	hashedPassword, err := password.HashPassword(regReq.Password)
	if err != nil {
		return RegisterResponse{}, errmsg.ErrorResponse{
			Message: err.Error(),
			Errors:  map[string]interface{}{"user_Register": err.Error()},
		}
	}
	var user = User{
		ID:           0,
		Username:     regReq.Username,
		FirstName:    regReq.FirstName,
		LastName:     regReq.LastName,
		PhoneNumber:  regReq.PhoneNumber,
		Email:        regReq.Email,
		Avatar:       regReq.Avatar,
		BirthDate:    regReq.BirthDate,
		IsActive:     true,
		Role:         0,
		PasswordHash: hashedPassword,
	}

	if _, err = srv.repository.RegisterUser(ctx, user); err != nil {
		return RegisterResponse{}, errmsg.ErrorResponse{
			Message: "registration failed",
			Errors:  map[string]interface{}{"user_Register": err.Error()},
		}
	}

	return RegisterResponse{
		ID: user.ID,
	}, nil
}
