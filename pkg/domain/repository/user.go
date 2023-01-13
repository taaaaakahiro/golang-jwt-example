package repository

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
)

type IUserRepository interface {
	GetUser(ctx context.Context, userID string) (*entity.User, error)
}
