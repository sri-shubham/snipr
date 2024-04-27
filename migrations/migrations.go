package migrations

import (
	"context"

	"github.com/sri-shubham/snipr/storage/persist/postgres"
	"github.com/uptrace/bun"
)

func MigrateDB(db *bun.DB) error {
	_, err := db.NewCreateTable().IfNotExists().Model(&postgres.PGShortenedURL{}).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
