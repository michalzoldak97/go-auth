package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db       *pgxpool.Pool
	security SecurityConfig
)

type Models struct {
	User  User
	Token Token
	SecurityConfig
}

func selectRows(query string, params ...any) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), security.DBTimeout)
	defer cancel()

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func selectRow(query string, params ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(context.Background(), security.DBTimeout)
	defer cancel()

	row := db.QueryRow(ctx, query, params...)

	return row
}

func New(dbPool *pgxpool.Pool) (Models, error) {

	db = dbPool

	return Models{
		User:           User{},
		Token:          Token{},
		SecurityConfig: SecurityConfig{},
	}, nil
}
