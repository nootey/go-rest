package workers

import (
	"context"
	"fmt"
	"go-rest/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func SeedAdminUser(ctx context.Context, client *mongo.Client, dbName string) error {
	collection := client.Database(dbName).Collection("users")

	creds, err := utils.LoadSeederCredentials()
	if err != nil {
		return fmt.Errorf("failed to load seeder credentials: %w", err)
	}
	adminEmail, ok := creds["ADMIN_EMAIL"]
	if !ok || adminEmail == "" {
		return fmt.Errorf("ADMIN_EMAIL not set in seeder credentials")
	}
	superAdminPassword, ok := creds["ADMIN_PASSWORD"]
	if !ok || superAdminPassword == "" {
		return fmt.Errorf("ADMIN_PASSWORD not set in seeder credentials")
	}

	hashedPassword, err := utils.HashAndSaltPassword(superAdminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := bson.M{
		"first_name":     "John",
		"last_name":      "Doe",
		"password":       hashedPassword,
		"email":          adminEmail,
		"email_verified": nil,
		"role":           "admin",
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
		"deleted_at":     nil,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	log.Println("Seeded admin user.")
	return nil
}
