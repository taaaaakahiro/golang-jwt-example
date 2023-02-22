package user

import (
	"encoding/json"
	"golang-jwt-example/pkg/domain/input"
	"golang-jwt-example/pkg/domain/output"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to read body", zap.Error(err))
			return
		}
		var in input.LoginInfo
		err = json.Unmarshal(body, &in)
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to unmarshal in", zap.Error(err))
			return
		}

		user, err := h.repo.UserRepository.GetUser(ctx, in.LoginID)
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to get user", zap.Error(err))
			return
		}
		if user == nil {
			msg := "unauthorized user"
			http.Error(w, msg, http.StatusUnauthorized)
			h.logger.Error("unauthorized user")
			return
		}

		// パスワード検証
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to compare hash and password", zap.Error(err))
			return
		}

		//JWT生成
		ulId := ulid.Make()
		now := time.Now().UTC()
		numericNow := jwt.NewNumericDate(now)
		accessTokenExpiredAt := now.Add(h.cfg.AccessTokenExpiredDuration)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ID:        ulId.String(),
			Subject:   in.LoginID,
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiredAt),
			NotBefore: numericNow,
			IssuedAt:  numericNow,
		})
		accessToken, err := token.SignedString([]byte(h.cfg.AccessTokenSecret))
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to sign access token", zap.Error(err))
			return
		}

		// トークンを保存(session方式でトークンを保持する場合)
		if err = h.redis.Conn.Set(ctx, in.LoginID, accessToken, 1*time.Hour).Err(); err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to set access token", zap.Error(err))
			return
		}

		b, err := json.Marshal(accessToken)
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to marshal access token", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
}
