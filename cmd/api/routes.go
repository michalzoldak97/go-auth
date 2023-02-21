package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res := `
			{
				"message": "Hello from backend"
			}
		`
		w.Write([]byte(res))
	})

	mux.Post("/v1/users/signup", app.signUp)
	mux.Post("/v1/users/login", app.login)

	mux.Get("/v1/users/validate", app.validateToken)

	return mux
}
