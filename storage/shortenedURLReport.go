//go:generate mockgen -source=shortenedURLReport.go -destination shortenReport_mock.go -package storage
package storage

import (
	context "context"

	models "github.com/sri-shubham/snipr/storage/models"
	"github.com/sri-shubham/snipr/storage/persist/postgres"
	"github.com/uptrace/bun"
)

type URLReport interface {
	ReportTopDomains(ctx context.Context, n int) ([]*models.JSONDomainReport, error)
}

func NewPGURLReport(db *bun.DB) URLReport {
	return &postgres.PGShortenedURLStorage{
		DB: db,
	}
}
