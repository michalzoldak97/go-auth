package main

import (
	"context"
	"log"

	"github.com/michalzoldak97/go-auth/env"
	"github.com/michalzoldak97/go-auth/pg"
)

func main() {
	err := env.LoadEnvVars("./.env")
	if err != nil {
		log.Fatalln("failed to load environment vars")
	}

	dsn, err := getDSN()
	if err != nil || dsn == "" {
		log.Fatalf("error while loading environment vars: %v\n", err)
	}

	dbPool, err := pg.NewPGPool(context.Background(), dsn)
	if err != nil {
		log.Fatalln("db connection failed")
	}
	defer dbPool.Close()

	app := &application{}
	app.loadApp()

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
