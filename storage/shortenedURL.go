//go:generate mockgen -source=shortenedURL.go -destination shorten_mock.go -package storage
package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
	rediscache "github.com/sri-shubham/snipr/storage/cache/redisCache"
	"github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/storage/persist/postgres"
	"github.com/uptrace/bun"
)

type URLStorage interface {
	StoreShortURL(ctx context.Context, shortUrl *models.ShortenedURL) error
	GetOriginalURL(ctx context.Context, shortURL string) (*models.ShortenedURL, error)
}

func NewPGShortenedURLStorage(db *bun.DB) URLStorage {
	return &postgres.PGShortenedURLStorage{
		DB: db,
	}
}

func NewRedisShortenedURLStorage(db *redis.Client) URLStorage {
	return rediscache.RedisShortenedURLStorage{
		Redis: db,
	}
}
