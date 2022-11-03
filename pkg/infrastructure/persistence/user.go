package persistence

import (
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/repository"
	"golang-jwt-example/pkg/io"
)

type UserRepository struct {
	database *io.SQLDatabase
}

var _ repository.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(db *io.SQLDatabase) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (u *UserRepository) GetUser(userID int) (*entity.User, error) {
	return nil, nil
}
