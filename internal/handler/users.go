package handler

import (
	"errors"
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
	err := helper.ReadJSONRequest(r, credential)
	if err != nil {
		h.logger.Error(
			"Failed while reading JSON request",
			"error", err,
		)

		if writeErr := helper.JSONError(w, http.StatusInternalServerError, "An internal server error occured"); writeErr != nil {
			h.logger.Error("failed to write JSON error response", "error", writeErr)
		}
		return
	}

	err = h.service.Register(ctx, *credential)
	if err != nil {
		if errors.Is(err, helper.ErrEmailAlreadyExists) {
			h.logger.Info(
				"User registration blocked: email already exists",
				"email", credential.Email,
			)

			if writeErr := helper.JSONError(w, http.StatusConflict, "User with this email already exists"); writeErr != nil {
				h.logger.Error("failed to write JSON error response", "error", writeErr)
			}
			return
		}

		h.logger.Error(
			"Failed while registering new user",
			"email", credential.Email,
			"error", err,
		)

		if writeErr := helper.JSONError(w, http.StatusInternalServerError, "An internal server error occured"); writeErr != nil {
			h.logger.Error("failed to write JSON error response", "error", writeErr)
		}
		return
	}

	if writeErr := helper.JSONResponse(w, http.StatusOK, "User registered successfully", ""); writeErr != nil {
		h.logger.Error("failed to write JSON success response", "error", writeErr)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	credential := &model.Credential{}

	err := helper.ReadJSONRequest(r, credential)
	if err != nil {
		h.logger.Error(
			"Failed while reading JSON request",
			"error", err,
		)

		if writeErr := helper.JSONError(w, http.StatusInternalServerError, "An internal server error occured"); writeErr != nil {
			h.logger.Error("failed to write JSON error response", "error", writeErr)
		}
		return
	}

	tokenStr, err := h.service.Login(ctx, *credential)
	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			h.logger.Info(
				"User login blocked: email does not exist",
				"email", credential.Email,
			)

			if writeErr := helper.JSONError(w, http.StatusNotFound, "User with this email does not exist"); writeErr != nil {
				h.logger.Error("failed to write JSON error response", "error", writeErr)
			}
			return
		}

		if errors.Is(err, helper.ErrWrongPassword) {
			h.logger.Info(
				"User login blocked: wrong password",
				"email", credential.Email,
			)

			if writeErr := helper.JSONError(w, http.StatusUnauthorized, "Wrong Password"); writeErr != nil {
				h.logger.Error("failed to write JSON error response", "error", writeErr)
			}
			return
		}

		h.logger.Error(
			"Failed while logging in user",
			"email", credential.Email,
			"error", err,
		)

		if writeErr := helper.JSONError(w, http.StatusInternalServerError, "An internal server error occured"); writeErr != nil {
			h.logger.Error("failed to write JSON error response", "error", writeErr)
		}
		return
	}

	if writeErr := helper.JSONResponse(w, http.StatusOK, "User login successfully", tokenStr); writeErr != nil {
		h.logger.Error("failed to write JSON success response", "error", writeErr)
	}
}
