package postgreSQL

import (
	"CalorieGuide-db/internal/lib/repeatable"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func New(ctx context.Context, maxAttempts int, timeout time.Duration) (pool *pgxpool.Pool, err error) {
	dsn := "postgresql://postgres:postgres@payment_postgres_tests:5432/postgres"
	const operation = "storage.postgreSQL.New"
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return fmt.Errorf("%s: %w", operation, err)
		}
		return nil
	}, maxAttempts, timeout*time.Second)
	return pool, nil
}
