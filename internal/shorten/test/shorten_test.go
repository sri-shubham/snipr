package test

import (
	"context"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sri-shubham/snipr/internal/shorten"
	"github.com/sri-shubham/snipr/storage"
	"github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/util"
	"github.com/stretchr/testify/require"
)

func TestShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageMock := storage.NewMockURLStorage(ctrl)

	shortener := shorten.NewShortener(
		4,
		6,
		8,
		"localhost:8080",
		storageMock,
	)

	longURL, err := url.Parse("https://en.wikipedia.org/wiki/URL_shortening")
	require.Nil(t, err)
	require.NotNil(t, longURL)

	expectedShortURL, err := url.Parse("https://localhost:8080/6H6EhC")
	require.Nil(t, err)
	require.NotNil(t, expectedShortURL)

	storageMock.EXPECT().GetOriginalURL(gomock.Any(), gomock.Any()).Return(nil, util.ErrNotFound)
	storageMock.EXPECT().StoreShortURL(gomock.Any(), &models.ShortenedURL{
		URL:      longURL,
		ShortURL: expectedShortURL,
	}).Return(nil)
	storageMock.EXPECT().GetOriginalURL(gomock.Any(), expectedShortURL.String()).Return(&models.ShortenedURL{
		URL:      longURL,
		ShortURL: expectedShortURL,
	}, nil)
	shortenedUrl, err := shortener.Shorten(context.Background(), longURL, 0)
	require.Nil(t, err)
	require.NotNil(t, shortenedUrl)

	// Mock returns already cached url
	storageMock.EXPECT().GetOriginalURL(gomock.Any(), gomock.Any()).Return(&models.ShortenedURL{
		URL: longURL,
	}, nil)
	shortenedUrl2, err := shortener.Shorten(context.Background(), longURL, 1000)
	require.Nil(t, err)
	require.NotNil(t, shortenedUrl2)

	// Shortening again should return same url
	require.Equal(t, shortenedUrl.URL.String(), shortenedUrl2.URL.String())
}

func TestShortenWhenMinLengthNotAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageMock := storage.NewMockURLStorage(ctrl)

	shortener := shorten.NewShortener(
		4,
		6,
		8,
		"localhost:8080",
		storageMock,
	)

	longURL, err := url.Parse("https://en.wikipedia.org/wiki/URL_shortening")
	require.Nil(t, err)
	require.NotNil(t, longURL)

	longURL2, err := url.Parse("https://en.wikipedia.org/wiki/URL_shortening2")
	require.Nil(t, err)
	require.NotNil(t, longURL2)

	expectedShortURL, err := url.Parse("https://localhost:8080/6H6EhC")
	require.Nil(t, err)
	require.NotNil(t, expectedShortURL)

	expectedShortURL2, err := url.Parse("https://localhost:8080/fUfhORo")
	require.Nil(t, err)
	require.NotNil(t, expectedShortURL)

	storageMock.EXPECT().GetOriginalURL(gomock.Any(), expectedShortURL.String()).Return(&models.ShortenedURL{
		URL: longURL2,
	}, nil)
	storageMock.EXPECT().GetOriginalURL(gomock.Any(), expectedShortURL2.String()).Return(nil, util.ErrNotFound)
	storageMock.EXPECT().StoreShortURL(gomock.Any(), &models.ShortenedURL{
		URL:          longURL,
		ShortURL:     expectedShortURL2,
		TTLInSeconds: 1000,
	}).Return(nil)
	storageMock.EXPECT().GetOriginalURL(gomock.Any(), expectedShortURL2.String()).Return(&models.ShortenedURL{
		URL:          longURL,
		ShortURL:     expectedShortURL2,
		TTLInSeconds: 1000,
	}, nil)
	shortenedUrl, err := shortener.Shorten(context.Background(), longURL, 1000)
	require.Nil(t, err)
	require.NotNil(t, shortenedUrl)
}
