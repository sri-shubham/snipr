package models

import (
	"net/url"
	"time"
)

type ShortenedURL struct {
	URL          *url.URL  `json:"url"`
	ShortURL     *url.URL  `json:"short_url"`
	TTLInSeconds int64     `json:"ttl_in_seconds,string"`
	CreatedAt    time.Time `json:"created_at"`
}

type JSONShortenedURL struct {
	URL          string    `json:"url"`
	ShortURL     string    `json:"short_url"`
	TTLInSeconds int64     `json:"ttl_in_seconds,string"`
	CreatedAt    time.Time `json:"created_at"`
}

type JSONDomainReport struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

func PresentJsonShortenedURLModel(in *ShortenedURL) *JSONShortenedURL {
	return &JSONShortenedURL{
		URL:          in.URL.String(),
		ShortURL:     in.ShortURL.String(),
		TTLInSeconds: in.TTLInSeconds,
		CreatedAt:    in.CreatedAt,
	}
}

func MapJsonShortenedURLModel(in *JSONShortenedURL) (*ShortenedURL, error) {
	origUrl, err := url.Parse(in.URL)
	if err != nil {
		return nil, err
	}

	shortUrl, err := url.Parse(in.ShortURL)
	if err != nil {
		return nil, err
	}

	return &ShortenedURL{
		URL:          origUrl,
		ShortURL:     shortUrl,
		TTLInSeconds: in.TTLInSeconds,
		CreatedAt:    in.CreatedAt,
	}, nil
}
