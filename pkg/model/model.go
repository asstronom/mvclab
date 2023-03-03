package model

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDB struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, connStr string) (*UserDB, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %s", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging db: %s", err)
	}
	return &UserDB{pool: pool}, nil
}

func (db *UserDB) Close() {
	db.pool.Close()
}
