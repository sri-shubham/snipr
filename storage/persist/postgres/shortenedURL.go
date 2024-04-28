package postgres

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/util"
	"github.com/uptrace/bun"
)

type PGShortenedURL struct {
	bun.BaseModel `bun:"table:short_url,alias:surl"`
	URL           string    `bun:"url"`
	Domain        string    `bun:"domain"`
	ShortURL      string    `bun:"short_url,pk"`
	Expires       time.Time `bun:"expires"`
	CreatedAt     time.Time `bun:"created_at"`
}

type PGShortenedURLDomainReport struct {
	bun.BaseModel `bun:"table:short_url,alias:surl"`
	Domain        string `bun:"domain"`
	Count         int    `json:"count"`
}

type PGShortenedURLStorage struct {
	DB *bun.DB
}

// GetOriginalURL implements storage.URLStorage.
func (p *PGShortenedURLStorage) GetOriginalURL(ctx context.Context, shortURL string) (*models.ShortenedURL, error) {
	pgShortenedURL := new(PGShortenedURL)
	err := p.DB.NewSelect().Model(pgShortenedURL).Where("short_url = ?", shortURL).Scan(ctx)
	if err != nil {
		return nil, util.PresentStorageErrors(err)
	}

	shortenedURL, err := presentPGShortenedURLModel(pgShortenedURL)
	if err != nil {
		return nil, util.PresentStorageErrors(err)
	}

	return shortenedURL, nil
}

// StoreShortURL implements storage.URLStorage.
func (p *PGShortenedURLStorage) StoreShortURL(ctx context.Context, shortenedURL *models.ShortenedURL) error {
	pgShortendedURL := mapPGShortenedURLModel(shortenedURL)
	pgShortendedURL.CreatedAt = time.Now()
	_, err := p.DB.NewInsert().Model(pgShortendedURL).
		On("Conflict (short_url) do nothing").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return util.PresentStorageErrors(err)
	}
	return nil
}

func (p *PGShortenedURLStorage) ReportTopDomains(ctx context.Context, n int) ([]*models.JSONDomainReport, error) {
	domains := []*PGShortenedURLDomainReport{}
	err := p.DB.NewSelect().Model(&domains).
		Group("domain").OrderExpr("count(1) desc").ColumnExpr("domain, count(1) count").
		Limit(n).Scan(ctx, &domains)
	if err != nil {
		return nil, util.PresentStorageErrors(err)
	}

	return presentPGShortenedURLReport(domains), nil
}

func mapPGShortenedURLModel(in *models.ShortenedURL) *PGShortenedURL {
	now := time.Now()
	expires := now.Add(time.Duration(in.TTLInSeconds) * time.Second)
	fmt.Println(in.URL.Host, in.URL.String())
	return &PGShortenedURL{
		Domain:   in.URL.Host,
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
		return nil, util.PresentStorageErrors(err)
	}

	shortUrl, err := url.Parse(in.ShortURL)
	if err != nil {
		return nil, util.PresentStorageErrors(err)
	}

	return &models.ShortenedURL{
		URL:          origUrl,
		ShortURL:     shortUrl,
		TTLInSeconds: int64(ttl),
		CreatedAt:    in.CreatedAt,
	}, nil
}

func presentPGShortenedURLReport(in []*PGShortenedURLDomainReport) []*models.JSONDomainReport {
	out := make([]*models.JSONDomainReport, 0, len(in))
	for _, item := range in {
		out = append(out, &models.JSONDomainReport{
			Domain: item.Domain,
			Count:  item.Count,
		})
	}
	return out
}
