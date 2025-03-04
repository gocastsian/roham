package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	errmsg "roham/pkg/err_msg"
	"roham/pkg/validator"
)

type ValidatorUserRepository interface {
}

type Validator struct {
	repo ValidatorUserRepository
}

func NewValidator(repo ValidatorUserRepository) Validator {
	return Validator{repo: repo}
}

func ReturnValidationError(err error) error {
	if errorsMap, ok := err.(validation.Errors); ok {
		return validator.NewError(errorsMap, validator.Nested, errmsg.ErrValidationFailed.Error())
	}
	return validator.NewError(err, validator.Nested, errmsg.ErrUnexpectedError.Error())
}
