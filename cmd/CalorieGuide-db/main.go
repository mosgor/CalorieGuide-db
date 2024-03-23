package main

import (
	"CalorieGuide-db/internal/config"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/storage/postgreSQL"
	"context"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting db-access", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgreSQL.New(context.TODO(), 3, cfg.Timeout)
	if err != nil {
		log.Error("failed to init storage", slg.Err(err))
		os.Exit(1)
	}
	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd, envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
