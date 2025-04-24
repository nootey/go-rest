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
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting database migrations")

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

	driver, err := mongodb.WithInstance(dbClient, &mongodb.Config{
		DatabaseName: cfg.DatabaseName,
	})
	if err != nil {
		logger.Fatal("Failed to create MongoDB driver", zap.Error(err))
	}

	migrationsDir := "./pkg/database/migrations"

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		logger.Fatal("Failed to initialize migrate instance", zap.Error(err))
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
			logger.Fatal("Failed to reset database", zap.Error(err))
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
			logger.Fatal("Invalid version format for force", zap.String("value", versionStr), zap.Error(convErr))
		}
		err = m.Force(version)
	default:
		logger.Fatal("Invalid migration type", zap.String("type", migrationType))
	}

	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal("Migration failed", zap.String("type", migrationType), zap.Error(err))
	}

	logger.Info("Migrations completed successfully", zap.String("type", migrationType))
}
