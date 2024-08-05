package di

import (
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/IkezawaYuki/c_root/infrastructre/driver"
	"github.com/go-redis/redis/v8"
)

var sessionRedisClient *redis.Client = nil
var redisClient *redis.Client = nil

func init() {
	if sessionRedisClient == nil {
		sessionRedisClient = redis.NewClient(&redis.Options{Addr: config.Env.SessionRedisHost})
		fmt.Println("session redis connect success")
	}
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{Addr: config.Env.RedisHost})
		fmt.Println("redis connect success")
	}
}

func NewSessionDriver() driver.RedisDriver {
	return driver.NewRedisDriver(sessionRedisClient)
}

func NewRedisDriver() driver.RedisDriver {
	return driver.NewRedisDriver(redisClient)
}
