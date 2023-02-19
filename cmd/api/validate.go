package main

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/michalzoldak97/go-auth/internal/data"
)

func (app *application) validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if !app.security.EmailDomainsRestricted {
		return true
	}

	domain := strings.Split(email, "@")[1]

	for _, allowed := range app.security.AllowedDomains {
		if domain == allowed {
			return true
		}
	}

	return false
}

func (app *application) validatePassASCII(phrase string) bool {

	pLen := len(phrase)

	if pLen < app.security.PassMinLen ||
		pLen > app.security.PassMaxLen {
		return false
	}

	isLower := !app.security.PassLower
	isUpper := !app.security.PassUpper
	isNum := !app.security.PassNum
	isSpec := !app.security.PassSpecial

	for _, char := range phrase {
		asciiVal := int(char)

		if !isLower && asciiVal > 96 && asciiVal < 123 {
			isLower = true
			continue
		}

		if !isUpper && asciiVal > 64 && asciiVal < 91 {
			isUpper = true
			continue
		}

		if !isNum && asciiVal > 47 && asciiVal < 58 {
			isNum = true
			continue
		}

		if !isSpec &&
			((asciiVal > 31 && asciiVal < 48) ||
				(asciiVal > 57 && asciiVal < 65) ||
				(asciiVal > 90 && asciiVal < 97) ||
				(asciiVal > 122 && asciiVal < 127)) {
			isSpec = true
			continue
		}

	}

	return isLower && isUpper && isNum && isSpec
}

func (app *application) validateNewUser(u data.User) error {
	if !app.validatePassASCII(u.Password) {
		return errors.New("password does not meet the minimum complexity requirements")
	}

	if !app.validateEmail(u.Email) {
		return errors.New("invalid email")
	}

	duplicate, err := app.models.User.GetByEmail(u.Email)
	if err != nil {
		return err
	}

	if len(duplicate) > 0 {
		return fmt.Errorf("user %v already exists", duplicate[0].Email)
	}

	return nil
}
