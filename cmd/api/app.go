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

type application struct {
	config   appConfig
	security data.SecurityConfig
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

	app.config.env = env
	app.config.port = port

	// load config from db

	app.security, err = app.models.SecurityConfig.GetConfig()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) loadApp(dbPool *pgxpool.Pool) error {
	var err error

	app.models, err = data.New(dbPool)
	if err != nil {
		return err
	}

	app.config = appConfig{}

	err = app.loadAppConfig()
	if err != nil {
		fmt.Println("Error loading config")
		return err
	}

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	erroroLog := log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.infoLog = infoLog
	app.errorLog = erroroLog

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
