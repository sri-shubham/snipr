package postgres

import (
	"sync"

	"github.com/sri-shubham/snipr/internal/config"
	"github.com/sri-shubham/snipr/util"
	"github.com/uptrace/bun"
)

var db *bun.DB
var once *sync.Once = &sync.Once{}

func GetDB(config *config.PostgresConfig) (*bun.DB, error) {
	var err error
	once.Do(func() {
		db, err = util.OpenPostgresConn(config)
	})
	if err != nil {

		return nil, err
	}
	return db, nil
}
