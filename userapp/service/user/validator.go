package user

import (
	"errors"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gocastsian/roham/pkg/statuscode"
	"regexp"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/validator"
)

var (
	ErrValidationRequiredAndNotZero = "phone or Social ID is required & must be greater than 0"
	ErrValidationPhoneDigitOnly     = "phone must contain only digits"
	ErrValidationPhoneLength        = "phone must be between 10 and 15 digits long"
	ErrUsernameLength               = "username length must be between 4 and 20 characters"
	ErrUsernameEmpty                = "username can not be empty"
	ErrUsernameFormat               = "username should not starts with numbers"
	ErrFirstNameEmpty               = "first name can not be empty"
	ErrFirstNameLength              = "first name must be between 4 and 20 characters"
	ErrLastNameEmpty                = "last name can not be empty"
	ErrLastNameLength               = "last name must be between 4 and 20 characters"
	ErrEmailEmpty                   = "email can not be empty"
	ErrPasswordEmpty                = "password can not be empty"
	ErrPasswordFormat               = "password should contain upper case letter, lower case letter and number"
	ErrUnvalidDate                  = "unvalid date"
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

func (v Validator) ValidateUsername(username string) error {

	err := validation.Validate(
		username,
		validation.Required.Error(ErrUsernameEmpty),
		validation.Match(regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{2,19}`)).Error(ErrUsernameFormat),
		validation.Length(4, 20).Error(ErrUsernameLength),
	)
	return err
}
func (v Validator) ValidateFirstName(firstname string) error {

	err := validation.Validate(
		firstname,
		validation.Required.Error(ErrFirstNameEmpty),

		validation.Length(4, 20).Error(ErrFirstNameLength),
	)
	return err
}
func (v Validator) ValidateLastName(lastname string) error {

	err := validation.Validate(
		lastname,
		validation.Required.Error(ErrLastNameEmpty),
		validation.Length(4, 20).Error(ErrLastNameLength),
	)
	return err
}
func (v Validator) ValidateEmail(email string) error {
	err := validation.Validate(
		email,
		validation.Required.Error(ErrEmailEmpty),
		is.Email,
	)
	return err
}
func (v Validator) ValidatePassword(password string) error {
	err := validation.Validate(
		password,
		validation.Length(8, 0),
		validation.Required.Error(ErrPasswordEmpty),
		validation.By(validateStrongPassword),
	)
	return err
}
func (v Validator) ValidateConfirmPassword(confirmPassword string, password string) error {
	err := validation.Validate(
		confirmPassword,
		validation.By(func(value interface{}) error {
			if confirmPassword != password {

				return errors.New("passwords don't match")
			}
			return nil
		}),
	)
	return err
}
func validateStrongPassword(value interface{}) error {
	s, _ := value.(string)

	var hasUpper, hasLower, hasDigit bool
	for _, c := range s {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

func (v Validator) ValidateRegistration(registerReq RegisterRequest) error {
	firstnameErr := v.ValidateFirstName(registerReq.FirstName)
	lastnameErr := v.ValidateLastName(registerReq.LastName)
	emailErr := v.ValidateEmail(registerReq.Email)
	usernameErr := v.ValidateUsername(registerReq.Username)
	passwordErr := v.ValidatePassword(registerReq.Password)

	ConfirmPasswordErr := v.ValidateConfirmPassword(registerReq.ConfirmPassword, registerReq.Password)

	errorsMap := make(map[string]interface{})

	if ConfirmPasswordErr != nil {
		errorsMap["confirm_password"] = ConfirmPasswordErr.Error()
	}
	if passwordErr != nil {
		errorsMap["password"] = passwordErr.Error()
	}
	if usernameErr != nil {
		errorsMap["username"] = usernameErr.Error()
	}
	if firstnameErr != nil {
		errorsMap["firstName"] = firstnameErr.Error()
	}
	if lastnameErr != nil {
		errorsMap["lastName"] = lastnameErr.Error()
	}
	if emailErr != nil {
		errorsMap["email"] = emailErr.Error()
	}
	if firstnameErr != nil || lastnameErr != nil || emailErr != nil || usernameErr != nil || passwordErr != nil || ConfirmPasswordErr != nil {
		return errmsg.ErrorResponse{
			Message:         "user validation has error",
			Errors:          errorsMap,
			InternalErrCode: statuscode.IntCodeUserValidation,
		}
	}
	return nil
}

func (v Validator) ValidateBirthDate(birthDate string) error {
	err := validation.Date("2006-01-02").Error(ErrUnvalidDate).Validate(birthDate)
	return err
}
