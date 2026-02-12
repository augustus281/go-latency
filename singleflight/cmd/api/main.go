package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/augustus281/go-latency/singleflight/internal/cache"
	"github.com/augustus281/go-latency/singleflight/internal/handler"
	"github.com/augustus281/go-latency/singleflight/internal/repository"
	"github.com/augustus281/go-latency/singleflight/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// PostgreSQL connection
	config, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	config.MaxConns = 20 // Max Connections
	config.MinConns = 5  // Min Connections

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Redis
	redisClient := cache.NewRedis(os.Getenv("REDIS_ADDR"))

	// Layers
	repo := repository.NewTemplateRepository(db)
	svc := service.NewTemplateService(redisClient, repo)
	h := handler.NewHandler(svc)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Server running on :8080")
	log.Fatal(server.ListenAndServe())
}
