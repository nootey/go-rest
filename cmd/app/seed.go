package main

import (
	"context"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		seedType := "help"
		if len(args) > 0 {
			seedType = args[0]
		}
		return runSeeders(seedType)
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

func runSeeders(seedType string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	if !isValidSeedType(seedType) {
		return fmt.Errorf("invalid seed type: %s", seedType)
	}

	cfg := config.LoadConfig()

	dbClient, err := database.ConnectToDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer func() {
		_ = dbClient.Disconnect(context.Background())
	}()

	ctx := context.Background()

	switch seedType {
	case "full", "basic":
		err = seeders.SeedDatabase(ctx, dbClient, seedType, cfg.DatabaseName)
		if err != nil {
			return fmt.Errorf("failed to seed database for type %s: %w", seedType, err)
		}
	case "help":
		// just print help? maybe log info?
		return nil
	default:
		return fmt.Errorf("unhandled seeder type: %s", seedType)
	}

	return nil
}
