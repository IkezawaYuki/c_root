package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Environment struct {
	CRootHost        string `envconfig:"C_ROOT_HOST" required:"true"`
	CRootPort        string `envconfig:"C_ROOT_PORT" required:"true"`
	CorsAllowOrigins string `envconfig:"CORS_ALLOW_ORIGINS"`

	ClientID     string `envconfig:"CLIENT_ID"`
	ClientSecret string `envconfig:"CLIENT_SECRET"`

	RedisHost        string `envconfig:"REDIS_HOST"`
	SessionRedisHost string `envconfig:"SESSION_REDIS_HOST"`
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
