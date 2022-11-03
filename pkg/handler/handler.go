package handler

import (
	"golang-jwt-example/pkg/handler/general"
	user "golang-jwt-example/pkg/handler/user"
	"golang-jwt-example/pkg/infrastructure/persistence"

	"time"

	"go.uber.org/zap"
)

type Handler struct {
	General *general.Handler
	User    *user.Handler
}

type Config struct {
	AccessTokenSecret           string
	RefreshTokenSecret          string
	AccessTokenExpiredDuration  time.Duration
	RefreshTokenExpiredDuration time.Duration
}

func NewHandler(logger *zap.Logger, repo *persistence.Repositories, cfg *Config) *Handler {
	userConfig := &user.Config{
		AccessTokenSecret:           cfg.AccessTokenSecret,
		RefreshTokenSecret:          cfg.RefreshTokenSecret,
		AccessTokenExpiredDuration:  cfg.AccessTokenExpiredDuration,
		RefreshTokenExpiredDuration: cfg.RefreshTokenExpiredDuration,
	}

	h := &Handler{
		General: general.NewHandler(logger),
		User:    user.NewHandler(logger.Named("user"), repo, userConfig),
	}

	return h
}
