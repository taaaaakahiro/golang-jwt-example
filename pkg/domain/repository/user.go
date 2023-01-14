package repository

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
)

type IUserRepository interface {
	ListUsers(ctx context.Context) ([]*entity.User, error)
}
