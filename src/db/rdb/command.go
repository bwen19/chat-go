package rdb

import (
	"context"
	"time"
)

func (c *Caches) SetCache(ctx context.Context, key string, val interface{}) error {
	return c.Set(ctx, key, val, 24*time.Hour).Err()
}

func (c *Caches) GetCache(ctx context.Context, key string, val interface{}) error {
	return c.Get(ctx, key).Scan(val)
}

func (c *Caches) DelCache(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

func (c *Caches) PushCacheList(ctx context.Context, key string, val interface{}) error {
	return c.RPush(ctx, key, val).Err()
}

func (c *Caches) GetCacheList(ctx context.Context, key string, val interface{}) error {
	return c.LRange(ctx, key, 0, -1).ScanSlice(val)
}

func (c *Caches) PopCacheList(ctx context.Context, key string, count int, val interface{}) error {
	length, err := c.LLen(ctx, key).Result()
	if err != nil {
		return err
	}

	num := int(length) - count
	if num <= 0 {
		return nil
	}
	return c.LPopCount(ctx, key, num).ScanSlice(val)
}

func (c *Caches) GetCacheKeys(ctx context.Context, pattern string) ([]string, error) {
	res := make([]string, 0)
	iter := c.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		res = append(res, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
