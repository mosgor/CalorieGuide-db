package main

import (
	"CalorieGuide-db/internal/config"
	food2 "CalorieGuide-db/internal/food/db"
	food "CalorieGuide-db/internal/food/handlers"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/storage/postgreSQL"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
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

	foodRepo := food2.NewRepository(storage, log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/products", food.NewFindAll(log, foodRepo))
	router.Post("/product", food.NewAdd(log, foodRepo))

	log.Info("starting server", slog.String("addr", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server: " + err.Error())
	}
	log.Error("server stopped")
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
