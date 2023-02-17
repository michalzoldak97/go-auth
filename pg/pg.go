package pg

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgInstance *pgxpool.Pool
	pgOnce     sync.Once
)

func NewPGPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, dsn)
		if err != nil {
			fmt.Println("error while connecting to the database")
			os.Exit(1)
		}

		pgInstance = db
	})

	err := pgInstance.Ping(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("*** Db connection success! ***")

	return pgInstance, nil
}
