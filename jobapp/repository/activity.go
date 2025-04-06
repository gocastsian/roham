package repository

import "context"

type Repository struct {
}

func New() Repository {
	return Repository{}
}

func (r Repository) SayHelloInPersian(ctx context.Context, name string) (string, error) {
	return "سلام " + name, nil
}
