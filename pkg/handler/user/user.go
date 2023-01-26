package user

import (
	"context"
	"encoding/json"
	"golang-jwt-example/pkg/infrastructure/persistence"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func (h *Handler) GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		id := vars["id"]
		user, err := h.repo.UserRepository.GetUser(ctx, id)
		if err != nil {
			msg := "failed to get user"
			h.logger.Error(msg, zap.Error(err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if user == nil {
			msg := "user is not found"
			h.logger.Error(msg, zap.Error(err))
			http.Error(w, msg, http.StatusNotFound)
			return
		}

		b, err := json.Marshal(user)
		if err != nil {
			msg := "failed to marshal error"
			h.logger.Error(msg, zap.Error(err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)

	})
}
