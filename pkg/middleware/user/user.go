package user

import (
	"go.uber.org/zap"
	"golang-jwt-example/pkg/infrastructure/persistence"
	"time"
)

type Middleware struct {
	logger *zap.Logger
	repo   *persistence.Repositories
	cfg    *Config
}

type Config struct {
	AccessTokenSecret           string
	RefreshTokenSecret          string
	AccessTokenExpiredDuration  time.Duration
	RefreshTokenExpiredDuration time.Duration
}

func NewMiddleware(logger *zap.Logger, repositories *persistence.Repositories, cfg *Config) *Middleware {
	return &Middleware{
		logger: logger,
		repo:   repositories,
		cfg:    cfg,
	}
}
