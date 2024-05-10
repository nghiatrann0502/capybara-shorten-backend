package redisClient

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	url := "redis://:redis_password@localhost:6379/0"
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	c := redis.NewClient(opts)

	if _, err = c.Ping(context.TODO()).Result(); err != nil {
		return nil, errors.New("could not connect to redis")
	}

	return c, nil

}
