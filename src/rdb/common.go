package rdb

import (
	"context"
	"time"
)

func (r *Redis) Set(ctx context.Context, key string, val interface{}) error {
	return r.clt.Set(ctx, key, val, 24*time.Hour).Err()
}

func (r *Redis) Get(ctx context.Context, key string, val interface{}) error {
	return r.clt.Get(ctx, key).Scan(val)
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.clt.Del(ctx, key).Err()
}

func (r *Redis) RPush(ctx context.Context, key string, val interface{}) error {
	return r.clt.RPush(ctx, key, val).Err()
}

func (r *Redis) LGetAll(ctx context.Context, key string, val interface{}) error {
	return r.clt.LRange(ctx, key, 0, -1).ScanSlice(val)
}

func (r *Redis) LPopCount(ctx context.Context, key string, count int, val interface{}) error {
	length, err := r.clt.LLen(ctx, key).Result()
	if err != nil {
		return err
	}

	num := int(length) - count
	if num <= 0 {
		return nil
	}
	return r.clt.LPopCount(ctx, key, num).ScanSlice(val)
}

func (r *Redis) Keys(ctx context.Context, pattern string) ([]string, error) {
	res := make([]string, 0)
	iter := r.clt.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		res = append(res, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
