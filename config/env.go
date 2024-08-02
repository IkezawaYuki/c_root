package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Environment struct {
	CorsAllowOrigins string `envconfig:"CORS_ALLOW_ORIGINS"`

	ClientID     string `envconfig:"CLIENT_ID"`
	ClientSecret string `envconfig:"CLIENT_SECRET"`
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
