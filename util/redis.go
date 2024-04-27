package util

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sri-shubham/snipr/internal/config"
)

func OpenRedisConn(conf *config.RedisConfig) (*redis.Client, error) {
	url := fmt.Sprintf("redis://%s:%s@%s:%d/%d", conf.User, conf.Password, conf.Host, conf.Port, conf.DB)
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
