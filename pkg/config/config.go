package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
	Port          string
	Release       bool
}

func LoadConfig() *Config {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	release, err := strconv.ParseBool(os.Getenv("RELEASE"))
	if err != nil {
		log.Fatalf("Invalid release mode: %v", err)
	}

	return &Config{
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDatabase: os.Getenv("MONGO_DB"),
		Port:          os.Getenv("PORT"),
		Release:       release,
	}
}
