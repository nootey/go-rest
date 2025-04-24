package main

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [type]",
	Short: "Run database migrations",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationType := "help"
		if len(args) > 0 {
			migrationType = args[0]
		}
		runMigrations(migrationType)
	},
}

func runMigrations(migrationType string) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	logger.Info("Starting database migrations")

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

	driver, err := mongodb.WithInstance(dbClient, &mongodb.Config{
		DatabaseName: cfg.DatabaseName,
	})
	if err != nil {
		log.Fatalf("Failed to create MongoDB driver: %v", err)
	}

	migrationsDir := "./pkg/database/migrations"

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate instance: %v", err)
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
			log.Fatalf("Failed to reset database: %v", err)
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
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			log.Fatalf("Invalid version for force: %v", err)
		}
		err = m.Force(version)
	default:
		log.Fatalf("Invalid migration type: %s", migrationType)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration error: %v", err)
	}

	logger.Info("Migrations completed successfully")
}
