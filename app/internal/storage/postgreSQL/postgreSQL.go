// Package postgreSQL provides a database connection pool with startup resilience.
package postgreSQL

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mosgor/CalorieGuide-db/internal/lib/repeatable"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// New creates a pgxpool.Pool using credentials from environment variables.
// It wraps connection initialization in a retry loop to survive transient network
// or database unavailability during application boot.
func New(ctx context.Context, maxAttempts int, timeout time.Duration) (pool *pgxpool.Pool, err error) {
	username := os.Getenv("username")
	password := os.Getenv("password")
	db := os.Getenv("db")
	address := os.Getenv("address")
	dsn := "postgresql://" + username + ":" + password + "@" + address + ":5432/" + db
	const operation = "storage.postgreSQL.New"
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return fmt.Errorf("%s: %w", operation, err)
		}
		return nil
	}, maxAttempts, timeout*time.Second)
	return pool, nil
}
