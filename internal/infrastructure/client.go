package infrastructure

import (
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMysqlConnection() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DatabaseUser,
		config.Env.DatabasePass,
		config.Env.DatabaseHost,
		config.Env.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func GetRedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Env.RedisAddr,
		DB:   0,
	})
}
