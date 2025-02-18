package layer

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
