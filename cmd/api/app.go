package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type appConfig struct {
	port         int
	maxPOSTBytes int64
}

type application struct {
	config   appConfig
	infoLog  *log.Logger
	errorLog *log.Logger
}

func (app *application) loadAppConfig() error {

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return err
	}

	max_post_bytes, err := strconv.Atoi(os.Getenv("MAX_POST_BYTES"))
	if err != nil {
		return err
	}

	app.config.port = port
	app.config.maxPOSTBytes = int64(max_post_bytes)

	return nil
}

func (app *application) loadApp() error {
	app.config = appConfig{}

	err := app.loadAppConfig()
	if err != nil {
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
