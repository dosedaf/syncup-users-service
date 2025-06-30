package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/dosedaf/syncup-users-service/helper"
	"github.com/dosedaf/syncup-users-service/internal/model"
	"github.com/dosedaf/syncup-users-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewUserHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	credential := &model.Credential{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err.Error())
		helper.JSONError(w, 500, "could not read body")
	}

	err = json.Unmarshal(body, credential)
	if err != nil {
		log.Print(err.Error())
		helper.JSONError(w, 500, "could not unmarshal")
	}

	err = h.service.Register(*credential)
	if err != nil {
		helper.JSONError(w, 500, "could not register")
	}

	helper.JSONResponse(w, 200, "")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	credential := &model.Credential{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err.Error())
		helper.JSONError(w, 500, "could not read body")
	}

	err = json.Unmarshal(body, credential)
	if err != nil {
		log.Print(err.Error())
		helper.JSONError(w, 500, "could not unmarshal")
	}

	tokenStr, err := h.service.Login(*credential)
	if err != nil {
		helper.JSONError(w, 500, "could not register")
	}

	helper.JSONResponse(w, 200, tokenStr)
}
