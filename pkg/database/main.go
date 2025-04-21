package database

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go-rest/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

var (
	mongoClient *mongo.Client
	once        sync.Once
)

// ConnectToDatabase initializes the MGM connection and returns the raw client.
func ConnectToDatabase(cfg *config.Config) (*mongo.Client, error) {
	var err error
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(cfg.MongoURI).
			SetConnectTimeout(10 * time.Second)

		mongoClient, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			return
		}

		err = mongoClient.Ping(context.Background(), nil)
		if err != nil {
			log.Printf("Failed to ping MongoDB: %v", err)
			return
		}

		err = mgm.SetDefaultConfig(nil, cfg.MongoDatabase, options.Client().ApplyURI(cfg.MongoURI))
		if err != nil {
			log.Printf("Failed to initialize MGM: %v", err)
			return
		}

		log.Println("Connected to MongoDB successfully.")
	})

	if err != nil {
		return nil, err
	}

	return mongoClient, nil
}

// DisconnectDatabase closes the MongoDB connection.
func DisconnectDatabase() error {
	if mongoClient == nil {
		return nil
	}
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		return fmt.Errorf("failed to disconnect MongoDB client: %w", err)
	}
	log.Println("Disconnected from MongoDB.")
	return nil
}

// PingDatabase checks MongoDB connectivity.
func PingDatabase() error {
	if mongoClient == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}
	return mongoClient.Ping(context.Background(), nil)
}
