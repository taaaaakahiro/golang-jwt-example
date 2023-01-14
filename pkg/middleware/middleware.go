package middleware

import (
	"golang-jwt-example/pkg/infrastructure/persistence"
	"time"

	"go.uber.org/zap"
)

type Middleware struct {
	User   *Middleware
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
	// studentConfig := &Config{
	// 	AccessTokenSecret:           cfg.AccessTokenSecret,
	// 	RefreshTokenSecret:          cfg.RefreshTokenSecret,
	// 	AccessTokenExpiredDuration:  cfg.AccessTokenExpiredDuration,
	// 	RefreshTokenExpiredDuration: cfg.RefreshTokenExpiredDuration,
	// }

	h := &Middleware{
		logger: logger,
		cfg:    cfg,
		// User:   NewMiddleware(logger.Named("user"), repo, studentConfig),
	}

	return h
}
