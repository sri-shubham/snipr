package storage

import (
	"context"

	"github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/storage/persist/postgres"
	"github.com/uptrace/bun"
)

type URLStorage interface {
	StoreShortURL(ctx context.Context, shortUrl *models.ShortenedURL) error
	GetOriginalURL(ctx context.Context, shortURL string) (*models.ShortenedURL, error)
}

func NewPGShortenedURLStorage(db *bun.DB) URLStorage {
	return postgres.PGShortenedURLStorage{
		DB: db,
	}
}
