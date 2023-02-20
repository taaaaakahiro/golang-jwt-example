package user

import (
	"encoding/json"
	"fmt"
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
		var loginInfo input.LoginInfo
		err = json.Unmarshal(body, &loginInfo)
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to unmarshal loginInfo", zap.Error(err))
			return
		}

		user, err := h.repo.UserRepository.GetUser(ctx, loginInfo.LoginID)
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
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
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
			Subject:   loginInfo.LoginID,
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

		redisClient := h.redis.Conn
		if err = redisClient.Set(ctx, "access_token", accessToken, 1*time.Hour).Err(); err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to set access token", zap.Error(err))
			return
		}

		fmt.Println(redisClient.Get(ctx, "access_token").String())

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
