package rediscache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/util"
)

type RedisShortenedURLStorage struct {
	Redis *redis.Client
}

// GetOriginalURL implements storage.URLStorage.
func (p RedisShortenedURLStorage) GetOriginalURL(ctx context.Context, shortURL string) (*models.ShortenedURL, error) {
	value, err := p.Redis.Get(ctx, shortURL).Result()
	if err != nil {
		return nil, util.PresentStorageErrors(err)
	}

	rShortenedURL := &models.JSONShortenedURL{}
	err = json.Unmarshal([]byte(value), rShortenedURL)
	if err != nil {
		return nil, util.PresentStorageErrors(err)
	}

	shortenedURL, err := models.MapJsonShortenedURLModel(rShortenedURL)
	if err != nil {
		return nil, err
	}

	return shortenedURL, nil
}

// StoreShortURL implements storage.URLStorage.
func (p RedisShortenedURLStorage) StoreShortURL(ctx context.Context, shortenedURL *models.ShortenedURL) error {
	redisShortenedURL := models.PresentJsonShortenedURLModel(shortenedURL)
	jsonBytes, err := json.Marshal(redisShortenedURL)
	if err != nil {
		return err
	}

	_, err = p.Redis.Set(ctx, shortenedURL.ShortURL.String(), string(jsonBytes),
		util.JitteredCacheDuration(util.DEFAULT_MIN_CACHE_TIME, util.DEFAULT_MAX_CACHE_TIME)).
		Result()
	if err != nil {
		return err
	}

	return nil
}
