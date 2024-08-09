package driver

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl int) error
	Del(ctx context.Context, key string) error
	GetClient() *redis.Client
}

func NewRedisDriver(client *redis.Client) RedisClient {
	return &redisClient{
		client: client,
	}
}

type redisClient struct {
	client *redis.Client
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

func (r *redisClient) Set(ctx context.Context, key string, value any, ttl int) error {
	_, err := r.client.Set(ctx, key, value, time.Second*time.Duration(ttl)).Result()
	return err
}

func (r *redisClient) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func (r *redisClient) GetClient() *redis.Client {
	return r.client
}
