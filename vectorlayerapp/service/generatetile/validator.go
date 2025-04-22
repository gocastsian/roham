package generatetile

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
