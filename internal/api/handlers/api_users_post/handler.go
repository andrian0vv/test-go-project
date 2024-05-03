package api_users_post

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/andrian0vv/test-go-project/internal/api"
	ie "github.com/andrian0vv/test-go-project/internal/errors"
	"github.com/andrian0vv/test-go-project/internal/services/users"
)

const (
	minAge            = 18
	minPasswordLength = 8
)

type userService interface {
	CreateUser(ctx context.Context, user users.User) error
}

type Handler struct {
	userService userService
	log         *slog.Logger
}

func New(userService userService, log *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		log:         log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req api.PostApiUsersJSONBody

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(api.ToError(ie.WithCode(err, api.CodeValidation)))
		_, _ = w.Write(body)
		h.log.Error("failed to decode body", err)

		return
	}

	if err := validate(req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		body, _ := json.Marshal(api.ToError(err))
		_, _ = w.Write(body)
		h.log.Error("invalid request", err)

		return
	}

	err := h.userService.CreateUser(r.Context(), convertUser(req))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		body, _ := json.Marshal(api.ToError(err))
		_, _ = w.Write(body)
		h.log.Error("failed to create user", err)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func validate(req api.PostApiUsersJSONBody) error {
	if req.Age < minAge {
		return ie.WithCode(fmt.Errorf("age should be more than %d", minAge), api.CodeValidation)
	}

	if len(req.Password) < minPasswordLength {
		return ie.WithCode(fmt.Errorf("password length should more than %d", minPasswordLength), api.CodeValidation)
	}

	return nil
}

func convertUser(req api.PostApiUsersJSONBody) users.User {
	return users.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		IsMarried: req.IsMarried,
		Password:  req.Password,
	}
}
