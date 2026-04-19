package main

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/jackc/pgx/v5/pgxpool"
    httpSwagger "github.com/swaggo/http-swagger"

    "ta-effective-mobile/internal/config"
    "ta-effective-mobile/internal/handler"
    "ta-effective-mobile/internal/repository"

    _ "ta-effective-mobile/docs"
)

// @title           Subscriptions Service API
// @version         1.0
// @description     REST-сервис для управления подписками пользователей
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
    log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    cfg, err := config.Load()
    if err != nil {
        log.Error("load config", "error", err)
        os.Exit(1)
    }

    db, err := pgxpool.New(context.Background(), cfg.DSN())
    if err != nil {
        log.Error("connect to database", "error", err)
        os.Exit(1)
    }
    defer db.Close()

    if err := db.Ping(context.Background()); err != nil {
        log.Error("ping database", "error", err)
        os.Exit(1)
    }
    log.Info("connected to database")

    m, err := migrate.New("file://migrations", cfg.DSN())
    if err != nil {
        log.Error("init migrations", "error", err)
        os.Exit(1)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Error("run migrations", "error", err)
        os.Exit(1)
    }
    log.Info("migrations applied")

    repo := repository.NewSubscriptionRepository(db)
    h := handler.NewSubscriptionHandler(repo, log)

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Route("/api/v1", func(r chi.Router) {
        r.Post("/subscriptions", h.Create)
        r.Get("/subscriptions", h.List)
        r.Get("/subscriptions/total", h.TotalCost)
        r.Get("/subscriptions/{id}", h.GetById)
        r.Put("/subscriptions/{id}", h.Update)
        r.Delete("/subscriptions/{id}", h.Delete)
    })

    r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.ServerPort)),
    ))

    addr := fmt.Sprintf(":%s", cfg.ServerPort)
    log.Info("server starting", "addr", addr)

    if err := http.ListenAndServe(addr, r); err != nil {
        log.Error("server error", "error", err)
        os.Exit(1)
    }
}