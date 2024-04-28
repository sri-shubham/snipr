package main

import (
	"log"
	"net/http"

	"github.com/sri-shubham/snipr/internal/config"
	"github.com/sri-shubham/snipr/internal/shorten"
	"github.com/sri-shubham/snipr/migrations"
	"github.com/sri-shubham/snipr/service"
	"github.com/sri-shubham/snipr/storage"
	rediscache "github.com/sri-shubham/snipr/storage/cache/redisCache"
	"github.com/sri-shubham/snipr/storage/persist/postgres"
)

func main() {
	log.Println("Loading config")
	config, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	log.Println("opening conn to db")
	pgDB, err := postgres.GetDB(config.Postgres)
	if err != nil {
		log.Fatalf("Failed to init postgres connection: %s", err)
	}

	log.Println("opening conn to redis")
	redis, err := rediscache.GetDB(config.Redis)
	if err != nil {
		log.Fatalf("Failed to init redis connection: %s", err)
	}

	log.Println("Running Migrations")
	err = migrations.MigrateDB(pgDB)
	if err != nil {
		log.Fatalf("Failed to init redis connection: %s", err)
	}

	_ = redis

	postgresURLStorage := storage.NewPGShortenedURLStorage(pgDB)
	postgresURLReport := storage.NewPGURLReport(pgDB)

	urlShorteningService := service.NewShortenURLService(
		shorten.NewShortener(
			config.Shortener.MinLength,
			config.Shortener.CustomMinLength,
			config.Shortener.CustomMaxLength,
			config.Host,
			postgresURLStorage,
		),
		postgresURLReport,
		postgresURLStorage,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", urlShorteningService.Shorten)
	mux.HandleFunc("POST /shorten/custom", urlShorteningService.ShortenCustom)
	mux.HandleFunc("GET /report/{count}", urlShorteningService.DomainReport)
	mux.HandleFunc("GET /{code}", urlShorteningService.Redirect)
	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
