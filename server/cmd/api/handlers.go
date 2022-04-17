package main

import (
	"net/http"
)


func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "Invalid json"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	// TODO: authenticate
	app.infoLog.Println(creds.Username, creds.Password)

	payload.Error = false
	payload.Message = "Authenticated"

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}
}
