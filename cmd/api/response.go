package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, app.security.MaxPOSTBytes)

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("duplicate JSON in the request body")
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusInternalServerError

	app.errorLog.Println(err)

	if len(status) > 0 {
		statusCode = status[0]
	}

	var res jsonResponse
	res.Error = true

	if app.config.env == "dev" ||
		err.Error() == "environment not specified" ||
		err.Error() == "password does not meet the minimum complexity requirements" {
		res.Message = fmt.Sprintf("error: %v", err)
	} else {
		res.Message = "something went wrong, contact the administrator"
	}

	app.writeJSON(w, statusCode, res)

	return nil
}
