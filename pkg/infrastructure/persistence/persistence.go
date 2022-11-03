package persistence

import "golang-jwt-example/pkg/io"

type Repositories struct {
	UserRepository *UserRepository
}

func NewRepositories(db *io.SQLDatabase) (*Repositories, error) {
	return &Repositories{
		UserRepository: NewUserRepository(db),
	}, nil
}
