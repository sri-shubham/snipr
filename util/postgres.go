package util

import (
	"database/sql"
	"fmt"

	"github.com/sri-shubham/snipr/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func OpenPostgresConn(conf *config.PostgresConfig) (*bun.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.DB)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	_, err := db.Exec("select 1;")
	if err != nil {
		return nil, err
	}

	return db, nil
}
