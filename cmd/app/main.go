package main

import (
	"context"
	"go-rest/internal/api/server"
	"go-rest/internal/services/mongo"
	"go-rest/pkg/config"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func main() {

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

	err = mongo.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("Unable to connect to database")
	}
	logger.Info("Connected to the database")

	// Initialize the server with the logger
	httpServer := server.NewServer(cfg, logger)

	// Start the server with health checks
	go httpServer.Start()

	// Wait for the interrupt signal
	<-ctx.Done()

	// Gracefully shutdown the HTTP server
	logger.Info("Shutting down server...")
	if err := httpServer.Shutdown(); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
