package infrastructure

import (
	"github.com/IkezawaYuki/popple/config"
	"github.com/redis/go-redis/v9"
)

func GetRedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Env.RedisAddr,
		DB:   0,
	})
}
