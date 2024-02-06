package app

import (
	"github.com/go-chi/chi/v5"
)

func (a *App) v1routes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Mount("/", a.routes())
	})

	return r
}

func (a *App) routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(a.handler.JSONHeader)

	// auth
	mux.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", a.handler.SignIn())
		})
	})

	// user
	mux.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", a.handler.CreateUser())
		})
		r.Route("/users/{user_id}", func(r chi.Router) {
			r.Get("/", a.handler.GetUser())
		})
	})

	return mux
}
