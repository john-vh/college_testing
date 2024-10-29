package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, keys ...string) error
}

var ErrNotFound = NotFoundError{}

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "Key not found"
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (cache *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return cache.client.Set(ctx, key, value, expiration).Err()
}

func (cache *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := cache.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNotFound
	}
	return val, nil
}

func (cache *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return cache.client.Del(ctx, keys...).Err()
}
