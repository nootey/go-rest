package main

import (
	"context"
	"github.com/spf13/cobra"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go-rest/pkg/database/seeders"
	"go.uber.org/zap"
	"log"
)

var seedCmd = &cobra.Command{
	Use:   "seed [type]",
	Short: "Run database seeders",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seedType := "help"
		if len(args) > 0 {
			seedType = args[0]
		}
		runSeeders(seedType)
	},
}

var validSeedTypes = map[string]bool{
	"full":  true,
	"basic": true,
	"help":  true,
}

func isValidSeedType(seedType string) bool {
	return validSeedTypes[seedType]
}

func runSeeders(seedType string) {

	// Validate seed type
	if !validSeedTypes[seedType] {
		log.Fatalf("Invalid seed type provided: %s.", seedType)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	logger.Info("Starting database seeding")

	cfg := config.LoadConfig()
	logger.Info("Loaded the configuration", zap.Any("config", cfg))

	dbClient, err := database.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error while disconnecting MongoDB client: %v", err)
		}
	}()

	ctx := context.Background()

	switch seedType {
	case "full":
		err = seeders.SeedDatabase(ctx, dbClient, "full", cfg.DatabaseName)
		if err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	case "basic":
		err = seeders.SeedDatabase(ctx, dbClient, "basic", cfg.DatabaseName)
		if err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	case "help":
		log.Fatal("\n Provide an additional argument to the seeder function. Valid arguments are: full, basic")
	default:
		log.Fatalf("Invalid seeder type: %s", seedType)
	}
}
