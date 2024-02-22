package app

import (
	"net/http"
)

func (a *App) v1Routes() *http.ServeMux {
	mux := http.NewServeMux()
	// mux.Use(a.handler.JSONHeader)

	// auth
	mux.HandleFunc("POST /v1/auth/signin", a.handler.SignIn())

	// user
	mux.HandleFunc("POST /v1/users", a.handler.CreateUser())
	mux.HandleFunc("GET /v1/users/{user_id}", a.handler.GetUser())

	return mux
}
