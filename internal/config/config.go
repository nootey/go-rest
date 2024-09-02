package config

import (
	"os"
)

type Config struct {
	MongoURI string
	Port     string
}

func LoadConfig() *Config {
	return &Config{
		MongoURI: os.Getenv("MONGO_URI"),
		Port:     os.Getenv("PORT"),
	}
}
