package platform

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustPGPool(ctx context.Context) *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URL") // e.g. postgres://user:pass@localhost:5432/app?sslmode=disable
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	return pool
}
