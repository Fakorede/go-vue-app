package main

import (
	"net/http"
	"server/internal/data"

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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/api/login", app.Login)
	mux.Post("/api/login", app.Login)

	mux.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
		var users data.User

		allUsers, err := users.GetAll()
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		app.writeJSON(w, http.StatusOK, allUsers)
	})

	mux.Get("/api/users/create", func(w http.ResponseWriter, r *http.Request) {
		var u = data.User{
			FirstName: "Anthon",
			LastName: "Freeze",
			Email: "anthon@test.com",
			Password: "password",
		}

		app.infoLog.Println("Creating user...")

		id, err := app.models.User.Insert(u)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			return
		}

		app.infoLog.Println("Created user id:", id)
		newUser, _ := app.models.User.GetByID(id)

		app.writeJSON(w, http.StatusOK, newUser)
	})

	return mux
}

