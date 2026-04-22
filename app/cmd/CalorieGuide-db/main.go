// Package main initializes the CalorieGuide API server.
// It loads configuration, establishes a database connection,
// configures HTTP middleware (CORS, JWT, logging, timeouts),
// registers route groups, and serves Swagger UI endpoints.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	client2 "github.com/mosgor/CalorieGuide-db/app/internal/client/db"
	client "github.com/mosgor/CalorieGuide-db/app/internal/client/handlers"
	"github.com/mosgor/CalorieGuide-db/app/internal/config"
	food2 "github.com/mosgor/CalorieGuide-db/app/internal/food/db"
	food "github.com/mosgor/CalorieGuide-db/app/internal/food/handlers"
	"github.com/mosgor/CalorieGuide-db/app/internal/lib/logger/slg"
	meal2 "github.com/mosgor/CalorieGuide-db/app/internal/meal/db"
	meal "github.com/mosgor/CalorieGuide-db/app/internal/meal/handlers"
	"github.com/mosgor/CalorieGuide-db/app/internal/storage/postgreSQL"

	_ "github.com/mosgor/CalorieGuide-db/app/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title CalorieGuide API
// @version 1.0
// @description API для работы с продуктами, пользователями и приёмами пищи
// @host localhost:8090
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting db-access", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgreSQL.New(context.TODO(), 3, 10)
	if err != nil {
		log.Error("failed to init storage", slg.Err(err))
		os.Exit(1)
	}
	if err := storage.Ping(context.Background()); err != nil {
		panic("Can't connect to postgres")
	}
	log.Info("Successfully connected to docker")

	foodRepo := food2.NewRepository(storage, log)
	clientRepo := client2.NewRepository(storage, log)
	mealRepo := meal2.NewRepository(storage, log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(5 * time.Second))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://158.160.212.231:80", "http://158.160.212.231"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Swagger UI
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://158.160.212.231:80/swagger/doc.json"), // указываем путь к swagger.json
	))

	router.Get("/redoc", func(w http.ResponseWriter, r *http.Request) {
		// Загружаем HTML из ReDoc
		html := `
<!DOCTYPE html>
<html>
  <head>
    <title>ReDoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='/swagger/doc.json'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
  </body>
</html>
`
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	})

	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.GetToken(log)))
		r.Use(jwtauth.Authenticator(config.GetToken(log)))

		// Product routes
		r.Post("/product", food.NewAdd(log, foodRepo))
		r.Put("/products/{id}", food.NewUpdate(log, foodRepo))
		r.Delete("/products/{id}", food.NewDelete(log, foodRepo))
		r.Post("/products/like", food.NewLike(log, foodRepo))

		// Client routes
		r.Put("/user/{id}", client.NewUpdate(log, clientRepo))
		r.Delete("/user/{id}", client.NewDelete(log, clientRepo, foodRepo, mealRepo))
		r.Get("/user/{id}/meals", client.NewFindMealLikes(log, clientRepo))
		r.Get("/user/{id}/products", client.NewFindFoodLikes(log, clientRepo))

		// Meal routes
		r.Post("/meal", meal.NewAdd(log, mealRepo, foodRepo))
		r.Post("/meals/like", meal.NewLike(log, mealRepo))
		r.Delete("/meals/{id}", meal.NewDelete(log, mealRepo))
		r.Put("/meals/{id}", meal.NewUpdate(log, mealRepo, foodRepo))
	})

	// Client routes
	router.Post("/user", client.NewAdd(log, clientRepo))
	router.Post("/login", client.FindEmail(log, clientRepo))

	// Product routes
	router.Post("/products", food.NewFindAll(log, foodRepo))
	router.Get("/products/{id}", food.NewFindOne(log, foodRepo))
	router.Post("/products/search", food.NewSearch(log, foodRepo))

	// Meal routes
	router.Post("/meals", meal.NewFindAll(log, mealRepo))
	router.Get("/meals/{id}", meal.NewFindOne(log, mealRepo))
	router.Post("/meals/search", meal.NewSearch(log, mealRepo))

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
