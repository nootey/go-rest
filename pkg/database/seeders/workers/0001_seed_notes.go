package workers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func SeedNotes(ctx context.Context, client *mongo.Client, dbName string) error {
	collection := client.Database(dbName).Collection("notes")

	note := bson.M{
		"title":       "My first note",
		"description": "Hello there",
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	_, err := collection.InsertOne(ctx, note)
	if err != nil {
		return err
	}

	log.Println("Seeded notes.")
	return nil
}
