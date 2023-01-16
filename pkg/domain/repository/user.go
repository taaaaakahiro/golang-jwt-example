package repository

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
)

type IUserRepository interface {
	ListUsers(ctx context.Context) ([]*entity.User, error)
	CreateUser(ctx context.Context, inputData input.User) (interface{}, error)
}
