package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/michalzoldak97/go-auth/internal/data"
)

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	var u data.User

	err := app.readJSON(w, r, &u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// email(username) -> trim whitespaces and convert to lowercase
	u.Email = strings.TrimSpace(u.Email)
	u.Email = strings.ToLower(u.Email)

	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	err = app.validateNewUser(u)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	newID, err := app.models.User.Create(u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "user created",
		Data:    envelope{"id": newID},
	}

	app.writeJSON(w, http.StatusCreated, res)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	if !app.validatePassASCII(req.Password) {
		err = errors.New("password does not meet the minimum complexity requirements")
		app.errorJSON(w, err)
		return
	}

	if !app.validateEmail(req.Email) {
		err = errors.New("invalid email")
		app.errorJSON(w, err)
		return
	}

	users, err := app.models.User.GetByEmail(req.Email)

	if err != nil || len(users) != 1 {
		err = errors.New("error fetching the user data")
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	user := users[0]

	if !user.IsPasswordValid(req.Password) {
		err = errors.New("invalid email/password")
		app.errorJSON(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "user found",
		Data:    envelope{"user": user},
	}

	app.writeJSON(w, http.StatusOK, res)
}
