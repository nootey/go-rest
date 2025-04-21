package main

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"go-rest/internal/bootstrap"
	"go-rest/internal/http"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go-rest/pkg/database/seeders"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:     "go-rest",
	Short:   "go-rest",
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

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

func init() {
	cobra.OnInitialize()

	// Register commands
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		zap.L().Fatal("Failed to execute the root command", zap.Error(err))
	}
}

func runServer() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logger.Sync()

	cfg := config.LoadConfig()
	logger.Info("Loaded the configuration", zap.Any("config", cfg))

	dbClient, err := database.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection Error: %v", err)
	}
	defer database.DisconnectDatabase()
	logger.Info("Connected to the database")

	// Initialize the server with the logger
	container := bootstrap.NewContainer(cfg, dbClient)
	httpLogger := logger.Named("http").With(zap.String("component", "HTTP"))
	httpServer := http.NewServer(container, httpLogger)
	go httpServer.Start()

	// Wait for the interrupt signal
	<-ctx.Done()

	// Gracefully shutdown the HTTP server
	logger.Info("Shutting down HTTP server...")
	if err := httpServer.Shutdown(); err != nil {
		logger.Fatal("HTTP Server forced to shutdown", zap.Error(err))
	}

	logger.Info("HTTP server exiting")
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

func runSeeders(seedType string) {

	// Validate seed type
	var validSeedTypes = map[string]bool{
		"full":  true,
		"basic": true,
		"help":  true,
	}

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
