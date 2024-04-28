package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sri-shubham/snipr/internal/shorten"
	"github.com/sri-shubham/snipr/storage"
	"github.com/sri-shubham/snipr/storage/models"
)

type ShortenUrlService interface {
	Shorten(w http.ResponseWriter, r *http.Request)
	ShortenCustom(w http.ResponseWriter, r *http.Request)
	DomainReport(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
}

type shortenURLServiceImpl struct {
	shortener shorten.Shortener
	report    storage.URLReport
	storage   storage.URLStorage
}

func NewShortenURLService(
	shortener shorten.Shortener,
	report storage.URLReport,
	storage storage.URLStorage,
) ShortenUrlService {
	return &shortenURLServiceImpl{
		shortener: shortener,
		report:    report,
		storage:   storage,
	}
}

type ShortenRequest struct {
	OriginalURL string    `json:"url"`
	Expires     time.Time `json:"expires"`
}

func (s *shortenURLServiceImpl) Shorten(w http.ResponseWriter, r *http.Request) {
	var requestBody ShortenRequest

	// Unmarshal the JSON data into the struct
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(requestBody.OriginalURL, "http://") && !strings.HasPrefix(requestBody.OriginalURL, "https://") {
		requestBody.OriginalURL = "https://" + requestBody.OriginalURL
	}

	requestUrl, err := url.Parse(requestBody.OriginalURL)
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to process request", http.StatusBadRequest)
		return
	}

	shortenedURL, err := s.shortener.Shorten(
		r.Context(),
		requestUrl,
		time.Duration(time.Until(requestBody.Expires)),
	)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, shorten.ErrNotAvailable) {
			code = http.StatusConflict
		}
		WriteJsonErrorResponseWithCode(w, err, "Failed to shorten url", code)
		return
	}

	resp := &models.JSONShortenedURL{
		URL:          shortenedURL.URL.String(),
		ShortURL:     shortenedURL.ShortURL.String(),
		TTLInSeconds: shortenedURL.TTLInSeconds,
		CreatedAt:    shortenedURL.CreatedAt,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	WriteJsonResponseWithCode(w, out, http.StatusOK)
}

type ShortenCustomRequest struct {
	OriginalURL string    `json:"url"`
	CustomCode  string    `json:"custom_code"`
	Expires     time.Time `json:"expires"`
}

func (s *shortenURLServiceImpl) ShortenCustom(w http.ResponseWriter, r *http.Request) {
	var requestBody ShortenCustomRequest

	// Unmarshal the JSON data into the struct
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(requestBody.OriginalURL, "http://") && !strings.HasPrefix(requestBody.OriginalURL, "https://") {
		requestBody.OriginalURL = "https://" + requestBody.OriginalURL
	}

	requestUrl, err := url.Parse(requestBody.OriginalURL)
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to process request", http.StatusBadRequest)
		return
	}

	shortenedURL, err := s.shortener.ShortenCustom(
		r.Context(),
		requestUrl,
		requestBody.CustomCode,
		time.Duration(time.Until(requestBody.Expires)),
	)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, shorten.ErrNotAvailable) {
			code = http.StatusConflict
		}
		WriteJsonErrorResponseWithCode(w, err, "Failed to shorten url", code)
		return
	}

	resp := &models.JSONShortenedURL{
		URL:          shortenedURL.URL.String(),
		ShortURL:     shortenedURL.ShortURL.String(),
		TTLInSeconds: shortenedURL.TTLInSeconds,
		CreatedAt:    shortenedURL.CreatedAt,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	WriteJsonResponseWithCode(w, out, http.StatusOK)
}

type ReportResponse struct {
	Items []*models.JSONDomainReport `json:"items"`
	Count int                        `json:"count"`
}

// DomainReport implements ShortenUrlService.
func (s *shortenURLServiceImpl) DomainReport(w http.ResponseWriter, r *http.Request) {
	count := r.PathValue("count")
	if count == "" {
		WriteJsonErrorResponseWithCode(w, errors.New("count not provided"), "Count is required", http.StatusBadRequest)
		return
	}

	countInt, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Count should be integer", http.StatusBadRequest)
		return
	}

	if countInt <= 0 {
		countInt = 5
	}

	reportItems, err := s.report.ReportTopDomains(r.Context(), int(countInt))
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to get domain report", http.StatusBadRequest)
		return
	}

	resp := ReportResponse{
		Items: reportItems,
		Count: len(reportItems),
	}
	out, err := json.Marshal(resp)
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	WriteJsonResponseWithCode(w, out, http.StatusOK)
}

// Redirect implements ShortenUrlService.
func (s *shortenURLServiceImpl) Redirect(w http.ResponseWriter, r *http.Request) {
	requestedURL := r.URL.String()

	shortURL, err := s.storage.GetOriginalURL(r.Context(), requestedURL)
	if err != nil {
		WriteJsonErrorResponseWithCode(w, err, "Failed to get domain report", http.StatusBadRequest)
		return
	}

	if shortURL.TTLInSeconds <= 0 {
		WriteJsonResponseWithCode(w, []byte("Not found"), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, shortURL.URL.String(), http.StatusFound)
}
