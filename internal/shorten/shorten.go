//go:generate mockgen -source=shorten.go -destination shorten_mock.go -package shorten
package shorten

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/jxskiss/base62"
	"github.com/sri-shubham/snipr/storage"
	"github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/util"
)

var ErrNotAvailable = errors.New("short url not available")

var customCodeRegexp *regexp.Regexp = regexp.MustCompile("[a-zA-Z1-9]+")

type Shortener interface {
	Shorten(ctx context.Context, url *url.URL, ttl time.Duration) (*models.ShortenedURL, error)
	ShortenCustom(ctx context.Context, url *url.URL, customString string, ttl time.Duration) (*models.ShortenedURL, error)
}

type shortenImpl struct {
	storage         storage.URLStorage
	minLength       int
	customMinLength int
	customMaxLength int
	host            string
}

func NewShortener(minLength int,
	customMinLength int,
	customMaxLength int,
	host string,
	storage storage.URLStorage) Shortener {
	return &shortenImpl{
		storage:         storage,
		minLength:       minLength,
		customMinLength: customMinLength,
		customMaxLength: customMaxLength,
		host:            host,
	}
}

func (s *shortenImpl) Shorten(ctx context.Context, url *url.URL, ttl time.Duration) (*models.ShortenedURL, error) {
	stringURL := url.String()
	hash := sha256.Sum256([]byte(stringURL))

	currentLen := s.minLength
	currentShortenUrl := ""
	for {
		shortCode := hash[:currentLen]
		encoded := base62.EncodeToString(shortCode)
		currentShortenUrl = fmt.Sprintf("https://%s/%s", s.host, encoded)

		existingUrl, err := s.storage.GetOriginalURL(ctx, currentShortenUrl)
		if err != nil {
			fmt.Println(err, errors.Is(err, util.ErrNotFound))
			if errors.Is(err, util.ErrNotFound) {
				break
			}
			return nil, err
		}
		if existingUrl != nil {
			if existingUrl.URL.String() == stringURL {
				// If this url is already shortended return existing one
				return existingUrl, nil
			}
			currentLen++

			if currentLen == len(hash) {
				return nil, ErrNotAvailable
			}
			continue
		}
	}

	shortUrl, err := url.Parse(currentShortenUrl)
	if err != nil {
		return nil, err
	}

	shortendUrl := &models.ShortenedURL{
		URL:          url,
		TTLInSeconds: int64(ttl),
		ShortURL:     shortUrl,
	}

	err = s.storage.StoreShortURL(ctx, shortendUrl)
	if err != nil {
		return nil, err
	}

	shortendUrl, err = s.storage.GetOriginalURL(ctx, currentShortenUrl)
	if err != nil {
		return nil, err
	}

	return shortendUrl, nil
}

// ShortenCustom implements Shortener.
func (s *shortenImpl) ShortenCustom(ctx context.Context, url *url.URL, customString string, ttl time.Duration) (*models.ShortenedURL, error) {
	if len(customString) < s.customMinLength && len(customString) > s.customMaxLength {
		return nil, fmt.Errorf("custom url code should be between %d, %d", s.customMinLength, s.customMaxLength)
	}

	if !customCodeRegexp.Match([]byte(customString)) {
		return nil, fmt.Errorf("custom url can only contain alphanumeric string")
	}

	currentShortenUrl := fmt.Sprintf("https://%s/%s", s.host, customString)
	existingURL, err := s.storage.GetOriginalURL(ctx, currentShortenUrl)
	if !errors.Is(err, util.ErrNotFound) {
		return nil, err
	}

	if existingURL.URL.String() == currentShortenUrl {
		return existingURL, nil
	} else if existingURL != nil {
		return nil, ErrNotAvailable
	}

	shortUrl, err := url.Parse(currentShortenUrl)
	if err != nil {
		return nil, err
	}

	shortendUrl := &models.ShortenedURL{
		URL:          url,
		TTLInSeconds: int64(ttl),
		ShortURL:     shortUrl,
	}

	err = s.storage.StoreShortURL(ctx, shortendUrl)
	if err != nil {
		return nil, err
	}

	shortendUrl, err = s.storage.GetOriginalURL(ctx, currentShortenUrl)
	if err != nil {
		return nil, err
	}

	return shortendUrl, nil
}
