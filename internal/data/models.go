package data

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db        *pgxpool.Pool
	dbTimeout time.Duration
	passCost  int
)

type Models struct {
	User  User
	Token Token
}

func selectRows(query string, params ...any) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func selectRow(query string, params ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	row := db.QueryRow(ctx, query, params...)

	return row
}

func loadEnvVars() error {
	dt, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	pc, err := strconv.Atoi(os.Getenv("PASS_COST"))

	if err != nil {
		return errors.New("failed to load model env vars")
	}

	dbTimeout = time.Second * time.Duration(dt)
	passCost = pc

	return nil
}

func New(dbPool *pgxpool.Pool) (Models, error) {

	err := loadEnvVars()
	if err != nil {
		return Models{}, err
	}

	db = dbPool

	return Models{
		User:  User{},
		Token: Token{},
	}, nil
}
