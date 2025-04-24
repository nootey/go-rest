package main

import (
	"context"
	"github.com/spf13/cobra"
	"go-rest/internal/bootstrap"
	"go-rest/internal/http"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the API server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func runServer() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	cfg := config.LoadConfig()
	logger.Info("Configuration loaded",
		zap.String("port", cfg.Port),
		zap.String("database", cfg.DatabaseName),
		zap.Bool("release", cfg.Release),
	)

	dbClient, err := database.ConnectToDatabase(cfg)
	if err != nil {
		logger.Fatal("Database connection failed", zap.Error(err))
	}
	defer database.DisconnectDatabase()
	logger.Info("Successfully connected to the database")

	// Initialize the server with the logger
	container := bootstrap.NewContainer(cfg, dbClient)
	httpLogger := logger.Named("http").With(zap.String("component", "HTTP"))
	httpServer := http.NewServer(container, httpLogger)
	go httpServer.Start()

	// Wait for the interrupt signal
	<-ctx.Done()

	// Gracefully shutdown the HTTP server
	logger.Info("Interrupt signal received, shutting down HTTP server...")
	if err := httpServer.Shutdown(); err != nil {
		logger.Fatal("HTTP server shutdown failed", zap.Error(err))
	}

	logger.Info("HTTP server exiting")
}
