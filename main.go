package main

import (
	"log"

	"github.com/sri-shubham/snipr/internal/config"
	"github.com/sri-shubham/snipr/migrations"
	rediscache "github.com/sri-shubham/snipr/storage/cache/redisCache"
	"github.com/sri-shubham/snipr/storage/persist/postgres"
)

func main() {
	config, err := config.ParseConfig("config/config_test.yml")
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	pgDB, err := postgres.GetDB(config.Postgres)
	if err != nil {
		log.Fatalf("Failed to init postgres connection: %s", err)
	}

	redis, err := rediscache.GetDB(config.Redis)
	if err != nil {
		log.Fatalf("Failed to init redis connection: %s", err)
	}

	err = migrations.MigrateDB(pgDB)
	if err != nil {
		log.Fatalf("Failed to init redis connection: %s", err)
	}

	_ = redis
}
