package main

import (
	client2 "CalorieGuide-db/internal/client/db"
	client "CalorieGuide-db/internal/client/handlers"
	"CalorieGuide-db/internal/config"
	food2 "CalorieGuide-db/internal/food/db"
	food "CalorieGuide-db/internal/food/handlers"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/storage/postgreSQL"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
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
	clientRepo := client2.NewRepository(storage, log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.GetToken(log)))
		r.Use(jwtauth.Authenticator(config.GetToken(log)))
		r.Post("/product", food.NewAdd(log, foodRepo))
		r.Put("/products/{id}", food.NewUpdate(log, foodRepo))
		r.Put("/user/{id}", client.NewUpdate(log, clientRepo))
	})

	router.Get("/products", food.NewFindAll(log, foodRepo))
	router.Post("/user", client.NewAdd(log, clientRepo))
	router.Post("/login", client.FindEmail(log, clientRepo))
	router.Get("/products/{id}", food.NewFindOne(log, foodRepo))

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
