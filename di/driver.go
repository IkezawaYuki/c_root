package di

import (
	"github.com/IkezawaYuki/c_root/config"
	"github.com/IkezawaYuki/c_root/internal/croot/infrastructre/driver"
	"github.com/go-redis/redis/v8"
	"log/slog"
)

var sessionRedisClient *redis.Client = nil

func init() {
	if sessionRedisClient == nil {
		sessionRedisClient = redis.NewClient(&redis.Options{Addr: config.Env.SessionRedisHost})
		slog.Info("session redis connect success")
	}
}

func NewSessionDriver() driver.RedisDriver {
	return driver.NewRedisDriver(sessionRedisClient)
}
