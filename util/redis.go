package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sri-shubham/snipr/internal/config"
)

func OpenRedisConn(conf *config.RedisConfig) (*redis.Client, error) {
	url := fmt.Sprintf("redis://%s:%s@%s:%d/%d", conf.User, conf.Password, conf.Host, conf.Port, conf.DB)
	var rdb *redis.Client
	count := 3
	for count > 0 {
		count--
		opt, err := redis.ParseURL(url)
		if err != nil {
			return nil, err
		}

		rdb := redis.NewClient(opt)
		_, err = rdb.Ping(context.Background()).Result()
		if err != nil && count == 0 {
			return nil, err
		} else if err == nil {
			break
		}
		log.Println("Failed try", count)
		time.Sleep(5 * time.Second)
	}

	return rdb, nil
}
