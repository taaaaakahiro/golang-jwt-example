package middleware

import (
	"context"
	"encoding/json"
	"golang-jwt-example/pkg/domain/output"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"go.uber.org/zap"
)

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			b, err := json.Marshal(output.NewHttpUnauthorized())
			if err != nil {
				http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
				m.logger.Error("failed to marshal output unauthorized", zap.Error(err))
				return
			}
			http.Error(w, string(b), http.StatusUnauthorized)
			return
		}

		accessToken := strings.Replace(authorization, "Bearer ", "", 1)
		ctx = SetCurrentAccessToken(ctx, &accessToken)

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.AccessTokenSecret), nil
		})
		if err != nil {
			b, err := json.Marshal(output.NewHttpUnauthorized())
			if err != nil {
				http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
				m.logger.Error("failed to marshal output unauthorized", zap.Error(err))
				return
			}
			http.Error(w, string(b), http.StatusUnauthorized)
			return
		}

		if claims.VerifyExpiresAt(time.Now().UTC().Unix(), false) == false {
			b, err := json.Marshal(output.NewHttpUnauthorized())
			if err != nil {
				http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
				m.logger.Error("failed to marshal output unauthorized", zap.Error(err))
				return
			}
			http.Error(w, string(b), http.StatusUnauthorized)
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			b, err := json.Marshal(output.NewHttpUnauthorized())
			if err != nil {
				http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
				m.logger.Error("failed to marshal output unauthorized", zap.Error(err))
				return
			}
			http.Error(w, string(b), http.StatusUnauthorized)
			return
		}
		userID := sub.(string)
		ctx = context.WithValue(ctx, "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
