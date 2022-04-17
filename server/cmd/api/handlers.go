package main

import (
	"errors"
	"net/http"
	"time"
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

	user, err := app.models.User.GetByEmail(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	validPass, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPass {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	// generate token
	token, err := app.models.Token.GenerateToken(user.ID, 24 * time.Hour)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.Token.Insert(*token, *user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error: false,
		Message: "Authenticated",
		Data: envelope {"token": token, "user": user},
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		Token string `json:"token"`
	}

	_ = app.readJSON(w, r, &reqPayload)

	_ = app.models.Token.Delete(reqPayload.Token)

	response := jsonResponse {
		Error: false,
		Message: "Logged out",
	}
	
	_ = app.writeJSON(w, http.StatusOK, response)
}
