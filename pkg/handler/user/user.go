package user

import (
	"context"
	"encoding/json"
	"golang-jwt-example/pkg/domain/output"
	"golang-jwt-example/pkg/infrastructure/persistence"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
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

func (h *Handler) LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ulid := ulid.Make()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ID:      ulid.String(),
			Subject: "subject",
		})
		accessToken, err := token.SignedString([]byte("access_token_secret"))
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to sign access token", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)

		b, _ := json.Marshal(accessToken)
		w.Write(b)
	})
}
