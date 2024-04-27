package postgres

import (
	"context"
	"net/url"
	"time"

	"github.com/sri-shubham/snipr/storage/models"
	"github.com/uptrace/bun"
)

type PGShortenedURL struct {
	bun.BaseModel `bun:"table:short_url,alias:surl"`
	URL           string    `bun:"url"`
	ShortURL      string    `bun:"short_url,pk"`
	Expires       time.Time `bun:"expires"`
	CreatedAt     time.Time `bun:"created_at"`
}

type PGShortenedURLStorage struct {
	DB *bun.DB
}

// GetOriginalURL implements storage.URLStorage.
func (p PGShortenedURLStorage) GetOriginalURL(ctx context.Context, shortURL string) (*models.ShortenedURL, error) {
	pgShortenedURL := new(PGShortenedURL)
	err := p.DB.NewSelect().Model(pgShortenedURL).Where("short_url = ?", shortURL).Scan(ctx)
	if err != nil {
		return nil, err
	}

	shortenedURL, err := presentPGShortenedURLModel(pgShortenedURL)
	if err != nil {
		return nil, err
	}

	return shortenedURL, nil
}

// StoreShortURL implements storage.URLStorage.
func (p PGShortenedURLStorage) StoreShortURL(ctx context.Context, shortenedURL *models.ShortenedURL) error {
	pgShortendedURL := mapPGShortenedURLModel(shortenedURL)
	pgShortendedURL.CreatedAt = time.Now()
	_, err := p.DB.NewInsert().Model(pgShortendedURL).
		On("Conflict (short_url) do nothing").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func mapPGShortenedURLModel(in *models.ShortenedURL) *PGShortenedURL {
	now := time.Now()
	expires := now.Add(time.Duration(in.TTLInSeconds) * time.Second)
	return &PGShortenedURL{
		URL:      in.URL.String(),
		ShortURL: in.ShortURL.String(),
		Expires:  expires,
	}
}

func presentPGShortenedURLModel(in *PGShortenedURL) (*models.ShortenedURL, error) {
	ttl := in.Expires.Sub(time.Now()) / time.Second
	if ttl < 0 {
		ttl = 0
	}

	origUrl, err := url.Parse(in.URL)
	if err != nil {
		return nil, err
	}

	shortUrl, err := url.Parse(in.ShortURL)
	if err != nil {
		return nil, err
	}

	return &models.ShortenedURL{
		URL:          origUrl,
		ShortURL:     shortUrl,
		TTLInSeconds: int64(ttl),
		CreatedAt:    in.CreatedAt,
	}, nil
}
