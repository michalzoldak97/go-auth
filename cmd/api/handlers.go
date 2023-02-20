package main

import (
	"errors"
	"fmt"
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

func (app *application) saveLoginAttempt(w http.ResponseWriter, ula data.UserLoginAttempt) {
	if !ula.Success {
		app.errorJSON(w, errors.New(ula.Message), http.StatusUnauthorized)
	}

	err := app.models.UserLoginAttempt.Create(ula)
	if err != nil {
		app.errorLog.Println(err)
	}
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

	if len(req.Email) > 254 {
		req.Email = "email too long"
	}

	if !app.validatePassASCII(req.Password) {
		app.saveLoginAttempt(w, data.UserLoginAttempt{
			Email:   req.Email,
			Message: "password does not meet the minimum complexity requirements",
			Success: false,
		})
		return
	}

	if !app.validateEmail(req.Email) {
		app.saveLoginAttempt(w, data.UserLoginAttempt{
			Email:   req.Email,
			Message: "invalid email",
			Success: false,
		})
		return
	}

	users, err := app.models.User.GetByEmail(req.Email)

	if err != nil || len(users) != 1 {
		app.saveLoginAttempt(w, data.UserLoginAttempt{
			Email:   req.Email,
			Message: "error fetching the user data",
			Success: false,
		})
		return
	}

	user := users[0]

	if !user.IsPasswordValid(req.Password) {
		app.saveLoginAttempt(w, data.UserLoginAttempt{
			Email:   req.Email,
			Message: "invalid email/password",
			Success: false,
		})
		return
	}

	token, err := app.models.Token.Create(user)
	if err != nil {
		app.saveLoginAttempt(w, data.UserLoginAttempt{
			Email:   req.Email,
			Message: fmt.Sprintf("token creation failed: %v", err),
			Success: false,
		})
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "logged in",
		Data:    envelope{"token": token},
	}

	app.writeJSON(w, http.StatusOK, res)

	app.saveLoginAttempt(w, data.UserLoginAttempt{
		Email:   req.Email,
		Message: "logged in",
		Success: true,
	})
}
