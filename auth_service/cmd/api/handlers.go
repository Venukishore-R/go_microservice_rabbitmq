package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Auth) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("password mismatch"), http.StatusUnauthorized)
		return
	}

	// log authentication
	if err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email)); err != nil {
		app.errorJSON(w, errors.New("error from logger service, while logging"+err.Error()))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user: %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, 200, payload)

}

func (app *Auth) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsondata, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service:1900/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsondata))
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")

	client := app.Client
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
