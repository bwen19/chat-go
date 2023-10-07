package rdb

import (
	"gochat/src/util"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	SessionCache
	UserCache
}

type Redis struct {
	rdb *redis.Client
}

func NewRedis(config *util.Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       0,
	})
	return &Redis{rdb}, nil
}
