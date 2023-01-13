package persistence

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "users"

type UserRepository struct {
	database *mongo.Collection
}

var _ repository.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		database: db.Collection(collection),
	}
}

func (u *UserRepository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	// u.database.Find()
	return nil, nil
}
