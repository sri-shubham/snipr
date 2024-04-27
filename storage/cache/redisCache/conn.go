package rediscache

import (
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/sri-shubham/snipr/internal/config"
	"github.com/sri-shubham/snipr/util"
)

var rDB *redis.Client
var once *sync.Once = &sync.Once{}

func GetDB(config *config.RedisConfig) (*redis.Client, error) {
	var err error
	once.Do(func() {
		rDB, err = util.OpenRedisConn(config)
	})
	if err != nil {
		return nil, err
	}

	return rDB, nil
}
