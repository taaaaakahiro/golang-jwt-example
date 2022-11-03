package user

import (
	"golang-jwt-example/pkg/infrastructure/persistence"
	"time"

	"go.uber.org/zap"
)

type Handler struct {
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

func NewHandler(logger *zap.Logger, repositories *persistence.Repositories, cfg *Config) *Handler {
	return &Handler{
		logger: logger,
		repo:   repositories,
		cfg:    cfg,
	}
}
