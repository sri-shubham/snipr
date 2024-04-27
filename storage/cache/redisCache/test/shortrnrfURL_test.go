package test

import (
	"context"
	"net/url"
	"testing"

	"github.com/sri-shubham/snipr/internal/config"
	rediscache "github.com/sri-shubham/snipr/storage/cache/redisCache"
	"github.com/sri-shubham/snipr/storage/models"
	"github.com/stretchr/testify/require"
)

var storage *rediscache.RedisShortenedURLStorage

func init() {
	config, err := config.ParseConfig("../../../../config/config_test.yml")
	if err != nil {
		panic(err)
	}

	rDB, err := rediscache.GetDB(config.Redis)
	if err != nil {
		panic(err)
	}

	storage = &rediscache.RedisShortenedURLStorage{
		Redis: rDB,
	}

	err = rDB.FlushDB(context.Background()).Err()
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
	require.Equal(t, returnedShortUrl.TTLInSeconds, shortendedURL.TTLInSeconds)
	require.NotNil(t, shortendedURL.CreatedAt)
}
