package user

import (
	"context"
	"encoding/json"
	"golang-jwt-example/pkg/infrastructure/persistence"
	"net/http"
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

func (h *Handler) ListHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		users, err := h.repo.UserRepository.ListUsers(context.Background())
		if err != nil {
			h.logger.Error("internal server error", zap.Error(err))
			return
		}
		b, err := json.Marshal(users)
		if err != nil {
			h.logger.Error("failed to unmarshal error", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
}
