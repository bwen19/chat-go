package rdb

import (
	"github.com/redis/go-redis/v9"
)

func New(client *redis.Client) *Caches {
	return &Caches{Client: client}
}

type Caches struct {
	*redis.Client
}
