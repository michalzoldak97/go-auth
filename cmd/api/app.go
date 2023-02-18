package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michalzoldak97/go-auth/internal/data"
)

type appConfig struct {
	env          string
	port         int
	maxPOSTBytes int64
}

type application struct {
	config   appConfig
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
}

func (app *application) loadAppConfig() error {

	// load config from env

	env := os.Getenv("ENV")
	if env == "" {
		return errors.New("environment not specified")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return err
	}

	max_post_bytes, err := strconv.Atoi(os.Getenv("MAX_POST_BYTES"))
	if err != nil {
		return err
	}

	app.config.env = env
	app.config.port = port
	app.config.maxPOSTBytes = int64(max_post_bytes)

	// load config from db

	return nil
}

func (app *application) loadApp(dbPool *pgxpool.Pool) error {
	app.config = appConfig{}

	err := app.loadAppConfig()
	if err != nil {
		return err
	}

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	erroroLog := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.infoLog = infoLog
	app.errorLog = erroroLog

	app.models, err = data.New(dbPool)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	fmt.Printf("auth api listens on %v\n", app.config.port)

	return srv.ListenAndServe()
}
