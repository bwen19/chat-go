package rdb

import "context"

type Cacher interface {
	SetCache(ctx context.Context, key string, val interface{}) error
	GetCache(ctx context.Context, key string, val interface{}) error
	DelCache(ctx context.Context, key string) error
	PushCacheList(ctx context.Context, key string, val interface{}) error
	GetCacheList(ctx context.Context, key string, val interface{}) error
	PopCacheList(ctx context.Context, key string, count int, val interface{}) error
	GetCacheKeys(ctx context.Context, pattern string) ([]string, error)
}

var _ Cacher = (*Caches)(nil)
