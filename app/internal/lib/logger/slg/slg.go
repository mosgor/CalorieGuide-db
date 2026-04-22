// Package slg provides structured logging helpers compatible with the standard log/slog package.
package slg

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// PgErr extracts PostgreSQL error metadata (Message, Detail, Where) and formats it
// as a single slog.Attr for consistent error logging.
func PgErr(err pgconn.PgError) slog.Attr {
	rt := err.Message + " " + err.Detail + " " + err.Where
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(rt),
	}
}
