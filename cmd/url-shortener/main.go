package main

import (
	"URL-shortener/internal/config"
	"URL-shortener/internal/http-server/handlers/delete"
	"URL-shortener/internal/http-server/handlers/redirect"
	"URL-shortener/internal/http-server/handlers/url/save"
	"URL-shortener/internal/http-server/middleware/logger"
	"URL-shortener/internal/lib/logger/handlers/slogpretty"
	storage "URL-shortener/internal/storage"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

const (
	envLocal = "LOCAL"
	envProd  = "PROD"
	envDev   = "DEV"
)

func main() {
	cfg := config.MustLoad()

	var log = setupLogger(strings.ToUpper(cfg.Env))
	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	db, err := sql.Open("postgres", cfg.DB_DSN)
	if err != nil {
		log.Error("Failed to open database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	st, err := storage.New(db)
	if err != nil {
		log.Error("Failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("Database initialized successfully")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("Url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, st))
		r.Delete("/{alias}", delete.New(log, st))
	})

	router.Get("/{alias}", redirect.New(log, st))

	log.Info("Starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Error("Server stopped", slog.String("address", cfg.HTTPServer.Address))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slogpretty.SetupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
