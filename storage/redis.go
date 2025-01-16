package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Db     int
	Addr   string
	client *redis.Client
}

var ctx = context.Background()

func NewRedisCache(db int, addr string) *RedisCache {
	cache := &RedisCache{db, addr, nil}
	cache.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})
	return cache
}

func (r *RedisCache) KeyExist(key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if exists > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *RedisCache) AddKey(key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
