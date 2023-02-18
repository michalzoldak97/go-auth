package main

import (
	"net/http"

	"github.com/michalzoldak97/go-auth/internal/data"
)

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	var u data.User

	err := app.readJSON(w, r, &u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.validateNewUser(u)
	if err != nil {
		app.errorJSON(w, err)
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
