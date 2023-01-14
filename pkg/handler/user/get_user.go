package user

import (
	"encoding/json"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/output"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// student, err := h.repo.UserRepository.GetUser("1")
		// if err != nil || student == nil {
		// 	http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
		// 	h.logger.Error("failed to get current student", zap.Error(err))
		// 	return
		// }

		outputUser := entity.User{
			ID:   1,
			Name: "",
		}

		b, err := json.Marshal(outputUser)
		if err != nil {
			http.Error(w, output.NewHttpInternalServerError(), http.StatusInternalServerError)
			h.logger.Error("failed to marshal output user", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
}
