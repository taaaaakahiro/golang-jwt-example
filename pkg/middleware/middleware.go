package middleware

import (
	"golang-jwt-example/pkg/infrastructure/persistence"
	"golang-jwt-example/pkg/middleware/user"
	"time"

	"go.uber.org/zap"
)

type Middleware struct {
	User   *user.Middleware
	logger *zap.Logger
	cfg    *Config
}

type Config struct {
	AccessTokenSecret           string
	RefreshTokenSecret          string
	AccessTokenExpiredDuration  time.Duration
	RefreshTokenExpiredDuration time.Duration
}

func NewMiddleware(logger *zap.Logger, repo *persistence.Repositories, cfg *Config) *Middleware {
	userConfig := &user.Config{
		AccessTokenSecret: cfg.AccessTokenSecret,
		//RefreshTokenSecret:          cfg.RefreshTokenSecret,
		AccessTokenExpiredDuration: cfg.AccessTokenExpiredDuration,
		//RefreshTokenExpiredDuration: cfg.RefreshTokenExpiredDuration,
	}

	h := &Middleware{
		logger: logger,
		cfg:    cfg,
		User:   user.NewMiddleware(logger.Named("user"), repo, userConfig),
	}

	return h
}
