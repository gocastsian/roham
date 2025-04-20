package user

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/validator"
)

var (
	ErrValidationRequiredAndNotZero = "phone or Social ID is required & must be greater than 0"
	ErrValidationPhoneDigitOnly     = "phone must contain only digits"
	ErrValidationPhoneLength        = "phone must be between 10 and 15 digits long"
	ErrAvatarSize                   = "avatar file size should be less than "
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

func (v Validator) ValidateAvatar(avatar Avatar, size int64, formats []string) error {
	return validation.Validate(
		avatar.FileHandler,
		validation.By(fileMimeTypeCheck(formats)),
		validation.By(fileSizeLimit(size*1024*1024)), // size MB
	)

}
func fileSizeLimit(maxBytes int64) validation.RuleFunc {
	return func(value interface{}) error {
		file, ok := value.(*multipart.FileHeader)
		if !ok || file == nil {
			return fmt.Errorf("invalid file size")
		}
		if file.Size > maxBytes {
			return fmt.Errorf(ErrAvatarSize, maxBytes)
		}
		return nil
	}
}
func fileMimeTypeCheck(allowedFormats []string) validation.RuleFunc {
	return func(value interface{}) error {
		fileHeader, ok := value.(*multipart.FileHeader)
		if !ok || fileHeader == nil {
			return fmt.Errorf("invalid file type")
		}
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		defer file.Close()

		// Read the first 512 bytes for content sniffing (as per net/http.DetectContentType)
		buf := make([]byte, 512)
		n, err := file.Read(buf)
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("could not read file: %w", err)
		}

		mimeType := http.DetectContentType(buf[:n])
		for _, allowedType := range allowedFormats {
			if mimeType == allowedType {
				return nil
			}
		}
		return fmt.Errorf("invalid mime type: %s; allowed: %v", mimeType, allowedFormats)
	}
}
