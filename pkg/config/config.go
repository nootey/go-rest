package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
	Port          string
	Release       string
}

func LoadConfig() *Config {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDatabase: os.Getenv("MONGO_DB"),
		Port:          os.Getenv("PORT"),
		Release:       os.Getenv("RELEASE"),
	}
}
