package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sri-shubham/snipr/internal/shorten"
	"github.com/sri-shubham/snipr/service"
	"github.com/sri-shubham/snipr/storage"
	"github.com/sri-shubham/snipr/storage/models"
	"github.com/stretchr/testify/require"
)

func TestShortenHTTPHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shortenMock := shorten.NewMockShortener(ctrl)

	shortenService := service.NewShortenURLService(shortenMock, nil)

	reqBody := &service.ShortenRequest{
		OriginalURL: "https://en.wikipedia.org/wiki/URL_shortening",
		Expires:     time.Time{}.Add(1000 * time.Second),
	}

	oURL, err := url.Parse(reqBody.OriginalURL)
	require.Nil(t, err)

	sURL, err := url.Parse("https://snipr.com/5rt3fv")
	require.Nil(t, err)

	bodyBytes, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(bodyBytes))
	respWriter := httptest.NewRecorder()

	shortenMock.EXPECT().Shorten(gomock.Any(), oURL, time.Until(reqBody.Expires)).Return(&models.ShortenedURL{
		URL:          oURL,
		ShortURL:     sURL,
		TTLInSeconds: 1000,
	}, nil)
	shortenService.Shorten(respWriter, req)
	require.Equal(t, respWriter.Result().StatusCode, http.StatusOK)

	resp := &models.JSONShortenedURL{}

	err = json.Unmarshal(respWriter.Body.Bytes(), &resp)
	require.Nil(t, err)

	require.Equal(t, resp.URL, oURL.String())
	require.Equal(t, resp.ShortURL, sURL.String())
	require.Equal(t, resp.TTLInSeconds, int64(1000))
}

func TestShortenHTTPHandlerInvalidURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shortenMock := shorten.NewMockShortener(ctrl)

	shortenService := service.NewShortenURLService(shortenMock, nil)

	reqBody := &service.ShortenRequest{
		OriginalURL: "https://en.wiki pedia.org/wiki/URL_shortening",
		Expires:     time.Time{}.Add(1000 * time.Second),
	}

	bodyBytes, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(bodyBytes))
	respWriter := httptest.NewRecorder()

	shortenService.Shorten(respWriter, req)
	require.Equal(t, respWriter.Result().StatusCode, http.StatusBadRequest)

	resp := &service.ErrorResponse{}

	err = json.Unmarshal(respWriter.Body.Bytes(), &resp)
	require.Nil(t, err)

	require.NotZero(t, resp.Error)
	require.NotZero(t, resp.Message)
}

func TestShortenCustomHTTPHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shortenMock := shorten.NewMockShortener(ctrl)

	shortenService := service.NewShortenURLService(shortenMock, nil)

	reqBody := &service.ShortenCustomRequest{
		OriginalURL: "https://en.wikipedia.org/wiki/URL_shortening",
		CustomCode:  "sniper",
		Expires:     time.Time{}.Add(1000 * time.Second),
	}

	oURL, err := url.Parse(reqBody.OriginalURL)
	require.Nil(t, err)

	sURL, err := url.Parse("https://snipr.com/sniper")
	require.Nil(t, err)

	bodyBytes, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(bodyBytes))
	respWriter := httptest.NewRecorder()

	shortenMock.EXPECT().ShortenCustom(gomock.Any(), oURL, "sniper", time.Until(reqBody.Expires)).Return(&models.ShortenedURL{
		URL:          oURL,
		ShortURL:     sURL,
		TTLInSeconds: 1000,
	}, nil)
	shortenService.ShortenCustom(respWriter, req)
	require.Equal(t, respWriter.Result().StatusCode, http.StatusOK)

	resp := &models.JSONShortenedURL{}

	err = json.Unmarshal(respWriter.Body.Bytes(), &resp)
	require.Nil(t, err)

	require.Equal(t, resp.URL, oURL.String())
	require.Equal(t, resp.ShortURL, sURL.String())
	require.Equal(t, resp.TTLInSeconds, int64(1000))
}

func TestShortenCustomHTTPHandlerWhenCodeIsAlreadyUsed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shortenMock := shorten.NewMockShortener(ctrl)

	shortenService := service.NewShortenURLService(shortenMock, nil)

	reqBody := &service.ShortenCustomRequest{
		OriginalURL: "https://en.wikipedia.org/wiki/URL_shortening",
		CustomCode:  "sniper",
		Expires:     time.Time{}.Add(1000 * time.Second),
	}

	oURL, err := url.Parse(reqBody.OriginalURL)
	require.Nil(t, err)

	bodyBytes, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(bodyBytes))
	respWriter := httptest.NewRecorder()

	shortenMock.EXPECT().ShortenCustom(gomock.Any(), oURL, "sniper", time.Until(reqBody.Expires)).Return(nil, shorten.ErrNotAvailable)
	shortenService.ShortenCustom(respWriter, req)
	require.Equal(t, respWriter.Result().StatusCode, http.StatusConflict)

	resp := &service.ErrorResponse{}

	err = json.Unmarshal(respWriter.Body.Bytes(), &resp)
	require.Nil(t, err)

	require.NotZero(t, resp.Error)
	require.NotZero(t, resp.Message)
}

func TestDomainReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := storage.NewMockURLReport(ctrl)
	shortenService := service.NewShortenURLService(nil, storage)

	req := httptest.NewRequest("GET", "/report/1", nil)
	req.SetPathValue("count", "5")
	respWriter := httptest.NewRecorder()

	storage.EXPECT().ReportTopDomains(gomock.Any(), 5).Return([]*models.JSONDomainReport{{}, {}, {}}, nil)
	shortenService.DomainReport(respWriter, req)
	require.Equal(t, respWriter.Result().StatusCode, http.StatusOK)

	resp := &service.ReportResponse{}

	err := json.Unmarshal(respWriter.Body.Bytes(), &resp)
	require.Nil(t, err)

	require.Equal(t, resp.Count, 3)
	require.Len(t, resp.Items, 3)
}
