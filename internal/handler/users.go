package handler

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"

	"github.com/pkg/errors"

	"github.com/andyfusniak/monolith/service"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

const (
	errCodeUserIDInvalid = "users/user-id-invalid"
)

type createUserRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *Handler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := log.WithContext(ctx)

		// request body
		req := createUserRequest{}
		if err := h.decode(w, r, &req); err != nil {
			cl.Warn("[app] createUserRequest body decode failed", err)
			clientError(w, http.StatusBadRequest, errCodeBadRequest, err.Error()) // 400
			return
		}
		message, ok := validateCreateUserRequest(&req)
		if !ok {
			cl.Warnf("[app] CreateUser: validation failed %q", message)
			clientError(w, http.StatusBadRequest, errCodeBadRequest, message) // 400
			return
		}

		// create a new user
		fmt.Printf("%#v\n", h)
		user, err := h.svc.CreateUser(ctx, *req.Email, *req.Password)
		if err != nil {
			cl.Errorf("[app] svc.CreateUser(ctx, req.Email=%q, req.Password=%q) unexpected error: %+v",
				*req.Email, "*****", err)
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}

		// successful response
		response := containerResponse{Data: user}
		h.respond(ctx, w, r, response, http.StatusCreated) // 201
		cl.Infof("[app] created new user %s with email=%s",
			user.ID, user.Email)
	}
}

func validateCreateUserRequest(req *createUserRequest) (string, bool) {
	// email
	if req.Email == nil {
		return "email attribute not set", false
	}
	_, err := mail.ParseAddress(*req.Email)
	if err != nil {
		return "email attribute must be a valid email address", false
	}

	// password
	if req.Password == nil {
		return "password request parameter not set", false
	}

	return "", true
}

func (h *Handler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := log.WithContext(ctx)

		userID := chi.URLParam(r, "user_id")
		if !isValidUserID(userID) {
			cl.Warnf("[auth] path parameter /users/%s invalid", userID)
			clientError(w, http.StatusUnprocessableEntity, errCodeUserIDInvalid,
				"user_id url path parameter is not a valid user id") // 422
			return
		}

		// create a new user
		user, err := h.svc.GetUser(ctx, userID)
		if err != nil {
			if err == service.ErrUserNotFound {
				cl.Infof("[app] user %q not found", userID)
				w.WriteHeader(http.StatusNotFound) // 404
				return
			}

			cl.Errorf("[app] svc.CreateUser(ctx, userID=%q) unexpected error: %+v",
				userID, err)
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}

		// successful response
		response := containerResponse{Data: user}
		h.respond(ctx, w, r, response, http.StatusCreated) // 201
		cl.Infof("[app] responding with user %q email %q", user.ID, user.Email)
	}
}

// auth

type signInRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *Handler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cl := log.WithContext(ctx)

		// request body
		req := signInRequest{}
		if err := h.decode(w, r, &req); err != nil {
			cl.Warn("[app] signInRequest body decode failed", err)
			clientError(w, http.StatusBadRequest, errCodeBadRequest, err.Error()) // 400
			return
		}
		message, ok := validateSignInRequest(&req)
		if !ok {
			cl.Warnf("[app] SignIn: validation failed %q", message)
			clientError(w, http.StatusBadRequest, errCodeBadRequest, message) // 400
			return
		}

		// signin user
		user, err := h.svc.VerifyUserPassword(ctx, *req.Email, *req.Password)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrUserWrongPassword) {
				cl.Infof("[app] verify user password failed for user email=%s", *req.Email)
				w.WriteHeader(http.StatusUnauthorized) // 401
				return
			}

			cl.Errorf("[app] svc.VerifyUserPassword(ctx, req.Email=%s, req.Password=*****) unexpected error: %+v", *req.Email, err)
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}

		// successful response
		response := containerResponse{Data: user}
		cl.Infof("[app] successful signin for user user_id=%s email=%s",
			user.ID, user.Email)
		h.respond(ctx, w, r, response, http.StatusCreated) // 201
	}
}

func validateSignInRequest(req *signInRequest) (string, bool) {
	// email
	if req.Email == nil {
		return "email attribute not set", false
	}
	_, err := mail.ParseAddress(*req.Email)
	if err != nil {
		return "email attribute must be a valid email address", false
	}

	// password
	if req.Password == nil {
		return "password request parameter not set", false
	}

	return "", true
}

var userIDExp = regexp.MustCompile(`^[A-HJ-NP-Za-km-z1-9]{22}$`)

func isValidUserID(userID string) bool {
	return userIDExp.Match([]byte(userID))
}
