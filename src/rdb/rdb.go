package rdb

import (
	"context"
	"gochat/src/util"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, val interface{}) error
	Get(ctx context.Context, key string, val interface{}) error
	Del(ctx context.Context, key string) error
	RPush(ctx context.Context, key string, val interface{}) error
	LGetAll(ctx context.Context, key string, val interface{}) error
	LPopCount(ctx context.Context, key string, count int, val interface{}) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Close()
}

type Redis struct {
	clt *redis.Client
}

func NewRedis(config *util.Config) (*Redis, error) {
	clt := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       0,
	})
	return &Redis{clt}, nil
}

func (r *Redis) Close() {
	r.clt.Close()
}
