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
		ulid := ulid.Make()
		now := time.Now().UTC()
		numericNow := jwt.NewNumericDate(now)
		accessTokenExpiredAt := now.Add(h.cfg.AccessTokenExpiredDuration)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ID:        ulid.String(),
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

// func (h *Handler) Login() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// inputStudentLogin := &input.StudentLogin{}
// dec := json.NewDecoder(r.Body)
// err := dec.Decode(&inputStudentLogin)
// if err != nil {
// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 	h.logger.Error("failed to decode input student login", zap.Error(err))
// 	return
// }

// unprocessableContent, err := inputStudentLogin.Validate()
// if err != nil {
// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 	h.logger.Error("failed to validate input student login", zap.Error(err))
// 	return
// }
// if unprocessableContent != nil {
// 	b, err := json.Marshal(unprocessableContent)
// 	if err != nil {
// 		http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 		h.logger.Error("failed to marshal output unprocessable content", zap.Error(err))
// 		return
// 	}
// 	http.Error(w, string(b), http.StatusUnprocessableEntity)
// 	return
// }

// user, err := h.repo.UserRepository.GetUser(inputStudentLogin.LoginID)
// if err != nil {
// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 	h.logger.Error("failed to get student", zap.Error(err))
// 	return
// }
// if user == nil {
// 	b, err := json.Marshal(output.NewHttpUnauthorized())
// 	if err != nil {
// 		http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 		h.logger.Error("failed to marshal http unauthorized", zap.Error(err))
// 		return
// 	}

// 	http.Error(w, string(b), http.StatusUnauthorized)
// 	return
// }

// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputStudentLogin.Password))
// if err != nil {
// 	b, err := json.Marshal(output.NewHttpUnauthorized())
// 	if err != nil {
// 		http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 		h.logger.Error("failed to marshal http unauthorized", zap.Error(err))
// 		return
// 	}

// 	http.Error(w, string(b), http.StatusUnauthorized)
// 	return
// }

// now := time.Now().UTC()
// numericNow := jwt.NewNumericDate(now)
// accessTokenExpiredAt := now.Add(h.cfg.AccessTokenExpiredDuration)

// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
// 	ID:        fmt.Sprintf("%d-%d", user.ID, now.UnixNano()),
// 	Subject:   user.Name,
// 	ExpiresAt: jwt.NewNumericDate(accessTokenExpiredAt),
// 	NotBefore: numericNow,
// 	IssuedAt:  numericNow,
// })
// accessToken, err := token.SignedString([]byte(h.cfg.AccessTokenSecret))
// if err != nil {
// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 	h.logger.Error("failed to sign access token", zap.Error(err))
// 	return
// }

// refreshTokenExpiredAt := now.Add(h.cfg.RefreshTokenExpiredDuration)
// token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
// 	ID:        fmt.Sprintf("%d-%d", user.ID, now.UnixNano()),
// 	Subject:   user.Name,
// 	ExpiresAt: jwt.NewNumericDate(refreshTokenExpiredAt),
// 	NotBefore: numericNow,
// 	IssuedAt:  numericNow,
// })
// refreshToken, err := token.SignedString([]byte(h.cfg.RefreshTokenSecret))
// if err != nil {
// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 	h.logger.Error("failed to sign refresh token", zap.Error(err))
// 	return
// }

// userToken := entity.StudentToken{
// 	StudentID:             student.ID,
// 	AccessToken:           accessToken,
// 	RefreshToken:          refreshToken,
// 	AccessTokenExpiredAt:  accessTokenExpiredAt,
// 	RefreshTokenExpiredAt: refreshTokenExpiredAt,
// }
// _, err = h.repo.StudentToken.CreateStudentToken(userToken, user.ID)
// if err != nil {
// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 	h.logger.Error("failed to create student token", zap.Error(err))
// 	return
// }

// outputLoginToken := output.LoginToken{
// 	AccessToken:  accessToken,
// 	RefreshToken: refreshToken,
// }

// b, err := json.Marshal(outputLoginToken)
// if err != nil {
// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
// 	h.logger.Error("failed to marshal output login token", zap.Error(err))
// 	return
// }

// 		w.WriteHeader(http.StatusCreated)
// 		// w.Write(b)
// 	})
// }
