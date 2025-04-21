package main

import (
	"context"
	"github.com/spf13/cobra"
	"go-rest/internal/api/server"
	"go-rest/internal/bootstrap"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go.uber.org/zap"
	"log"
	"os/signal"
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
func init() {
	cobra.OnInitialize()
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
	httpServer := server.NewServer(container, httpLogger)
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
