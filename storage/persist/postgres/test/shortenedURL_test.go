package test

import (
	"context"
	"net/url"
	"testing"

	"github.com/sri-shubham/snipr/internal/config"
	"github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/storage/persist/postgres"
	"github.com/stretchr/testify/require"
)

var storage *postgres.PGShortenedURLStorage

func init() {
	config, err := config.ParseConfig("../../../../config/config_test.yml")
	if err != nil {
		panic(err)
	}

	pgDB, err := postgres.GetDB(config.Postgres)
	if err != nil {
		panic(err)
	}

	storage = &postgres.PGShortenedURLStorage{
		DB: pgDB,
	}

	_, err = pgDB.NewTruncateTable().Model(&postgres.PGShortenedURL{}).Exec(context.Background())
	if err != nil {
		panic(err)
	}
}

func TestStoreShortURL(t *testing.T) {
	origUrl, _ := url.Parse("github.com/sri-shubham/Snipr")
	shortUrl, _ := url.Parse("snipr.com/shubham")
	shortendedURL := &models.ShortenedURL{
		URL:          origUrl,
		ShortURL:     shortUrl,
		TTLInSeconds: 10000,
	}

	err := storage.StoreShortURL(context.Background(), shortendedURL)
	require.Nil(t, err)

	returnedShortUrl, err := storage.GetOriginalURL(context.Background(), shortUrl.String())
	require.Nil(t, err)
	require.Equal(t, shortendedURL.ShortURL.String(), returnedShortUrl.ShortURL.String())
	require.Equal(t, shortendedURL.URL.String(), returnedShortUrl.URL.String())
	require.LessOrEqual(t, returnedShortUrl.TTLInSeconds, shortendedURL.TTLInSeconds)
	require.NotNil(t, shortendedURL.CreatedAt)
}
