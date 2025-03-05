package user

import (
	"context"
	"log/slog"

	errmsg "roham/pkg/err_msg"
)

type Repository interface {
	GetAllUsers(ctx context.Context) ([]User, error)
}

type Service struct {
	repository Repository
	validator  Validator
	logger     *slog.Logger
}

func NewService(repo Repository, validator Validator, logger *slog.Logger) Service {
	return Service{
		repository: repo,
		validator:  validator,
		logger:     logger,
	}
}

func (srv Service) GetAllUsers(ctx context.Context) (GetAllUserResponse, error) {

	users, err := srv.repository.GetAllUsers(ctx)
	if err != nil {
		srv.logger.Error("user_GetAllUsers", err)
		return GetAllUserResponse{}, errmsg.ErrorResponse{
			Message: err.Error(),
			Errors: map[string]interface{}{
				"user_GetAllUsers": err.Error(),
			},
		}
	}
	responseUsers := make([]GetUserResponse, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, GetUserResponse{
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
		})
	}
	return GetAllUserResponse{Users: responseUsers}, nil
}
