package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Environment struct {
	RedisAddr    string `envconfig:"REDIS_ADDR"`
	RedisPass    string `envconfig:"REDIS_PASS"`
	DatabaseUser string `envconfig:"DATABASE_USER"`
	DatabasePass string `envconfig:"DATABASE_PASS"`
	DatabaseName string `envconfig:"DATABASE_NAME"`
	DatabaseHost string `envconfig:"DATABASE_HOST"`
}

var Env Environment

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := envconfig.Process("", &Env); err != nil {
		log.Fatal(err)
	}
}
