package models

import (
	"net/url"
	"time"
)

type ShortenedURL struct {
	URL          url.URL   `json:"url"`
	ShortURL     url.URL   `json:"short_url"`
	TTLInSeconds int64     `json:"ttl_in_seconds"`
	CreatedAt    time.Time `json:"created_at"`
}
