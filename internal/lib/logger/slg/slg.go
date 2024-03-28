package slg

import (
	"github.com/jackc/pgconn"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func PgErr(err pgconn.PgError) slog.Attr {
	rt := err.Message + " " + err.Detail + " " + err.Where
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(rt),
	}
}
