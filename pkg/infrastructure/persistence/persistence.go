package persistence

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	UserRepository *UserRepository
}

func NewRepositories(db *mongo.Database) (*Repositories, error) {
	return &Repositories{
		UserRepository: NewUserRepository(db),
	}, nil
}
