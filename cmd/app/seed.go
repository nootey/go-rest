package main

import (
	"context"
	"github.com/spf13/cobra"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go-rest/pkg/database/seeders"
	"go.uber.org/zap"
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

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error()) // can't log if logger failed
	}
	defer logger.Sync()

	// Validate seed type
	if !isValidSeedType(seedType) {
		logger.Fatal("Invalid seed type provided", zap.String("seedType", seedType))
	}

	logger.Info("Starting database seeding")

	cfg := config.LoadConfig()
	logger.Info("Loaded the configuration", zap.Any("config", cfg))

	dbClient, err := database.ConnectToDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			logger.Error("Error while disconnecting MongoDB client", zap.Error(err))
		}
	}()

	ctx := context.Background()

	switch seedType {
	case "full", "basic":
		err = seeders.SeedDatabase(ctx, dbClient, seedType, cfg.DatabaseName)
		if err != nil {
			logger.Fatal("Failed to seed database", zap.String("type", seedType), zap.Error(err))
		}
		logger.Info("Database seeding completed for type: %s", zap.String("type", seedType))
	case "help":
		logger.Fatal("Seeder usage help", zap.Strings("validTypes", []string{"full", "basic"}))
	default:
		logger.Fatal("Unhandled seeder type", zap.String("type", seedType))
	}

}
