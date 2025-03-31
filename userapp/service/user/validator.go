package user

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/validator"
)

var (
	ErrValidationRequiredAndNotZero = "phone or Social ID is required & must be greater than 0"
	ErrValidationPhoneDigitOnly     = "phone must contain only digits"
	ErrValidationPhoneLength        = "phone must be between 10 and 15 digits long"
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

func (v Validator) ValidatePhoneNumber(phoneNumber string) error {
	err := validation.Validate(
		phoneNumber,
		validation.Required.Error(ErrValidationRequiredAndNotZero),
		validation.Match(regexp.MustCompile("^[0-9]+$")).Error(ErrValidationPhoneDigitOnly),
		validation.Length(10, 15).Error(ErrValidationPhoneLength),
	)
	return err
}
