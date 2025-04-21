package database

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go-rest/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
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

		err = mgm.SetDefaultConfig(nil, cfg.DatabaseName, options.Client().ApplyURI(cfg.MongoURI))
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

func ResetDatabase(client *mongo.Client, dbName string) error {
	ctx := context.Background()
	db := client.Database(dbName)

	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	for _, coll := range collections {
		err := db.Collection(coll).Drop(ctx)
		if err != nil {
			return fmt.Errorf("failed to drop collection %s: %w", coll, err)
		}
		log.Printf("Dropped collection: %s", coll)
	}

	return nil
}

func findMigrationFilesByVersion(migrationsDir string, version uint) ([]string, error) {
	var matched []string
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("%d", version)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			matched = append(matched, file.Name())
		}
	}

	return matched, nil
}

func LogMigrationFiles(logger *zap.Logger, dir string, version uint) {
	files, err := findMigrationFilesByVersion(dir, version)
	if err != nil {
		logger.Warn("Could not find migration files", zap.Error(err))
		return
	}
	for _, f := range files {
		logger.Info("Executed migration file", zap.String("file", f))
	}
}
