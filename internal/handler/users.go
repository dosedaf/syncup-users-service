package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/dosedaf/syncup-users-service/helper"
	"github.com/dosedaf/syncup-users-service/internal/model"
	"github.com/dosedaf/syncup-users-service/internal/service"
)

type HandlerInstance interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	service service.ServiceInstance
	logger  *slog.Logger
}

func NewUserHandler(service service.ServiceInstance, logger *slog.Logger) HandlerInstance {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	credential := &model.Credential{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Info("failed to register", "error", err)
		helper.JSONError(w, 500, "failed to register")
		return
	}

	err = json.Unmarshal(body, credential)
	if err != nil {
		h.logger.Info("failed to register", "error", err)
		helper.JSONError(w, 500, "failed to register")
		return
	}

	err = h.service.Register(ctx, *credential)
	if err != nil {
		h.logger.Info("failed to register", "error", err)
		helper.JSONError(w, 500, "failed to register")
		return
	}

	helper.JSONResponse(w, 200, "")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	credential := &model.Credential{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Info("failed to login", "error", err)
		helper.JSONError(w, 500, "failed to login")
		return
	}

	err = json.Unmarshal(body, credential)
	if err != nil {
		h.logger.Info("failed to login", "error", err)
		helper.JSONError(w, 500, "failed to login")
		return
	}

	tokenStr, err := h.service.Login(ctx, *credential)
	if err != nil {
		h.logger.Info("failed to login", "error", err)
		helper.JSONError(w, 500, "failed to login")
		return
	}

	helper.JSONResponse(w, 200, tokenStr)
}
