package main

import (
	"fmt"

	"github.com/michalzoldak97/go-auth/internal/data"
)

func (app *application) validateNewUser(u data.User) error {
	duplicate, err := app.models.User.GetByEmail(u.Email)
	if err != nil {
		return err
	}

	if len(duplicate) > 0 {
		return fmt.Errorf("user %v already exists", duplicate[0].Email)
	}

	return nil
}
