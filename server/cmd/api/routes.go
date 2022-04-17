package main

import (
	"net/http"
	"server/internal/data"
	"time"

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

	mux.Post("/api/login", app.Login)
	mux.Post("/api/logout", app.Logout)

	mux.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
		var users data.User

		allUsers, err := users.GetAll()
		if err != nil {
			app.errorLog.Println(err)
			return
		}


		response := jsonResponse{
			Error: false,
			Message: "success",
			Data: envelope{"users": allUsers},
		}

		app.writeJSON(w, http.StatusOK, response)
	})

	mux.Get("/api/users/create", func(w http.ResponseWriter, r *http.Request) {
		var u = data.User{
			FirstName: "Anthon",
			LastName:  "Freeze",
			Email:     "anthon@test.com",
			Password:  "password",
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

	mux.Get("/api/test-generate-token", func(w http.ResponseWriter, r *http.Request) {
		token, err := app.models.User.Token.GenerateToken(1, 60*time.Minute)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		token.Email = "admin@test.com"
		token.CreatedAt = time.Now()
		token.UpdatedAt = time.Now()

		payload := jsonResponse{
			Error:   false,
			Message: "success",
			Data:    token,
		}

		app.writeJSON(w, http.StatusOK, payload)
	})

	mux.Get("/api/test-save-token", func(w http.ResponseWriter, r *http.Request) {
		token, err := app.models.User.Token.GenerateToken(2, 60*time.Minute)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		user, err := app.models.User.GetByID(2)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		token.UserID = user.ID
		token.CreatedAt = time.Now()
		token.UpdatedAt = time.Now()

		err = token.Insert(*token, *user)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		payload := jsonResponse{
			Error:   false,
			Message: "success",
			Data:    token,
		}

		app.writeJSON(w, http.StatusOK, payload)
	})

	mux.Get("/api/test-valid-token", func(w http.ResponseWriter, r *http.Request) {
		tokenToValidate := r.URL.Query().Get("token")
		isValid, err := app.models.Token.ValidateToken(tokenToValidate)
		if err != nil {
			app.errorJSON(w, err)
		}

		var payload jsonResponse
		payload.Error = !isValid
		payload.Data = isValid

		app.writeJSON(w, http.StatusOK, payload)
	})

	return mux
}
