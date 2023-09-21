package rdb

import "github.com/redis/go-redis/v9"

type RedisDb interface {
}
type Redis struct {
	rdb *redis.Client
}

func NewRedis(addr string) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &Redis{rdb}, nil
}
