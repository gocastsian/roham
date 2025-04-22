package importlayer

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrValidationRequiredLess100Char = "is required & less than 100 characters"
)

type ValidatorRepository interface {
}

type Validator struct {
	repo ValidatorRepository
}

func NewValidator(repo ValidatorRepository) Validator {
	return Validator{
		repo: repo,
	}
}

func (v Validator) ValidateJobRequest(job CreateJobRequest) error {
	return validation.ValidateStruct(&job,
		validation.Field(&job.Token,
			validation.Required.Error(ErrValidationRequiredLess100Char),
			validation.Length(1, 100),
		),
	)
}
