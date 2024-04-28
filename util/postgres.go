package util

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/sri-shubham/snipr/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func OpenPostgresConn(conf *config.PostgresConfig) (*bun.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.DB)

	var db *bun.DB
	count := 3
	for count > 0 {
		count--

		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		db = bun.NewDB(sqldb, pgdialect.New())
		_, err := db.Exec("select 1;")
		if err != nil && count == 0 {
			return nil, err
		} else if err == nil {
			break
		}
		log.Println("Failed try", count)
		time.Sleep(time.Second * 5)
	}

	// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db, nil
}
