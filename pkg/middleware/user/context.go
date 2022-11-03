package user

import (
	"context"
	"errors"
	"golang-jwt-example/pkg/domain/entity"
)

type contextKey string

const currentAccessTokenKey = "currentAccessToken"
const currentStudentKey contextKey = "currentStudent"

func SetCurrentAccessToken(ctx context.Context, currentAccessToken *string) context.Context {
	return context.WithValue(ctx, currentAccessTokenKey, currentAccessToken)
}

func GetCurrentAccessToken(ctx context.Context) (*string, error) {
	v := ctx.Value(currentAccessTokenKey)
	accessToken, ok := v.(*string)
	if !ok {
		return nil, errors.New("current access token not found")
	}
	return accessToken, nil
}

func SetCurrentStudent(ctx context.Context, currentStudent *entity.User) context.Context {
	return context.WithValue(ctx, currentStudentKey, currentStudent)
}

func GetCurrentStudent(ctx context.Context) (*entity.User, error) {
	v := ctx.Value(currentStudentKey)
	student, ok := v.(*entity.User)
	if !ok {
		return nil, errors.New("current student not found")
	}
	return student, nil
}
