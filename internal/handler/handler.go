package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/andyfusniak/monolith/service"
	log "github.com/sirupsen/logrus"
)

const (
	// General
	errCodeBadRequest = "errors/bad-request"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

type containerResponse struct {
	Data any `json:"data"`
}

// apiErrorResponse standard response format for RESTful API calls.
type apiErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *Handler) decode(w http.ResponseWriter, r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		return err
	}
	return nil
}

func (h *Handler) respond(ctx context.Context, w http.ResponseWriter, r *http.Request, data any, status int) {
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			cl := log.WithContext(ctx)
			cl.Errorf("[app] response failed to json encode data")
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
	}
}

// 4xx (Client Error): The request contains bad syntax or cannot be fulfilled
func clientError(w http.ResponseWriter, statusCode int, code string, message string) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(apiErrorResponse{
		statusCode,
		code,
		message,
	})
}
