package driver

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisDriver interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl int) error
	Del(ctx context.Context, key string) error
	GetClient() *redis.Client
}

func NewRedisDriver(client *redis.Client) RedisDriver {
	return &redisDriver{
		client: client,
	}
}

type redisDriver struct {
	client *redis.Client
}

func (r *redisDriver) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

func (r *redisDriver) Set(ctx context.Context, key string, value any, ttl int) error {
	_, err := r.client.Set(ctx, key, value, time.Second*time.Duration(ttl)).Result()
	return err
}

func (r *redisDriver) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func (r *redisDriver) GetClient() *redis.Client {
	return r.client
}
