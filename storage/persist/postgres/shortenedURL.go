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
	URL           url.URL   `bun:"url"`
	ShortURL      url.URL   `bun:"short_url,pk"`
	Expires       time.Time `bun:"expires"`
	CreatedAt     time.Time `bun:"created_at"`
}

func (u *PGShortenedURL) BeforeInsert(ctx context.Context, query *bun.InsertQuery) error {
	u.CreatedAt = time.Now()
	return nil
}

type PGShortenedURLStorage struct {
	DB *bun.DB
}

// GetOriginalURL implements storage.URLStorage.
func (p PGShortenedURLStorage) GetOriginalURL(ctx context.Context, shortURL string) (*models.ShortenedURL, error) {
	shortenedURL := new(PGShortenedURL)
	err := p.DB.NewSelect().Model(shortenedURL).Where("short_url = ?", shortURL).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return presentPGShortenedURLModel(shortenedURL), nil
}

// StoreShortURL implements storage.URLStorage.
func (p PGShortenedURLStorage) StoreShortURL(ctx context.Context, shortenedURL *models.ShortenedURL) error {
	pgShortendedURL := mapPGShortenedURLModel(shortenedURL)
	_, err := p.DB.NewInsert().Model(pgShortendedURL).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func mapPGShortenedURLModel(in *models.ShortenedURL) *PGShortenedURL {
	Expires := time.Now().Add(time.Duration(in.TTLInSeconds) * time.Second)
	return &PGShortenedURL{
		URL:      in.URL,
		ShortURL: in.ShortURL,
		Expires:  Expires,
	}
}

func presentPGShortenedURLModel(in *PGShortenedURL) *models.ShortenedURL {
	ttl := time.Since(in.Expires) / time.Second
	if ttl < 0 {
		ttl = 0
	}

	return &models.ShortenedURL{
		URL:          in.URL,
		ShortURL:     in.ShortURL,
		TTLInSeconds: int64(ttl),
		CreatedAt:    in.CreatedAt,
	}
}
