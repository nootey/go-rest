package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	MongoURI           string
	DatabaseName       string
	Port               string
	Release            bool
	WebClientDomain    string
	WebClientPort      string
	JwtWebClientAccess string
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
		MongoURI:           os.Getenv("MONGO_URI"),
		DatabaseName:       os.Getenv("DATABASE_NAME"),
		Port:               os.Getenv("PORT"),
		WebClientDomain:    os.Getenv("WEB_CLIENT_DOMAIN"),
		WebClientPort:      os.Getenv("WEB_CLIENT_PORT"),
		JwtWebClientAccess: os.Getenv("JWT_WEB_CLIENT_ACCESS"),
		Release:            release,
	}
}
