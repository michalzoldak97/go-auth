package pg

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgPool struct {
	db *pgxpool.Pool
}

var (
	pgInstance *pgPool
	pgOnce     sync.Once
)

func (p *pgPool) Close() {
	p.db.Close()
}

func (p *pgPool) Ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}

func NewPGPool(ctx context.Context, dsn string) (*pgPool, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, dsn)
		if err != nil {
			fmt.Println("error while connecting to the database")
			os.Exit(1)
		}

		pgInstance = &pgPool{db}
	})

	err := pgInstance.db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("*** Db connection success! ***")

	return pgInstance, nil
}
