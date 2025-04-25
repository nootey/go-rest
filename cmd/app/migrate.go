package main

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go.uber.org/zap"
	"os"
	"strconv"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [type]",
	Short: "Run database migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		migrationType := "help"
		if len(args) > 0 {
			migrationType = args[0]
		}
		return runMigrations(migrationType)
	},
}

func runMigrations(migrationType string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	logger.Info("Starting database migrations")

	cfg := config.LoadConfig()
	logger.Info("Loaded the configuration", zap.Any("config", cfg))

	dbClient, err := database.ConnectToDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			logger.Error("Error while disconnecting MongoDB client", zap.Error(err))
		}
	}()

	driver, err := mongodb.WithInstance(dbClient, &mongodb.Config{
		DatabaseName: cfg.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("failed to create MongoDB driver: %w", err)
	}
	migrationsDir := "./pkg/database/migrations"

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	// Run based on the type
	var newVersion uint
	switch migrationType {
	case "up":
		err = m.Up()
		if err == nil || err == migrate.ErrNoChange {
			newVersion, _, _ = m.Version()
			database.LogMigrationFiles(logger, migrationsDir, newVersion)
		}
	case "down":
		err = m.Steps(-1)
	case "fresh":
		err = database.ResetDatabase(dbClient, cfg.DatabaseName)
		if err != nil {
			return fmt.Errorf("failed to reset database: %w", err)
		}
		err = m.Up()
		if err == nil || err == migrate.ErrNoChange {
			newVersion, _, _ = m.Version()
			database.LogMigrationFiles(logger, migrationsDir, newVersion)
		}
	case "reset":
		err = m.Down()
	case "force":
		versionStr := os.Getenv("VERSION")
		version, convErr := strconv.Atoi(versionStr)
		if convErr != nil {
			return fmt.Errorf("invalid version format for force: %w", convErr)
		}
		err = m.Force(version)
	default:
		return fmt.Errorf("invalid migration type: %s", migrationType)
	}

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed for type %s: %w", migrationType, err)
	}

	logger.Info("Migrations completed successfully", zap.String("type", migrationType))
	return nil
}
