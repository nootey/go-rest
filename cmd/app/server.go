package main

import (
	"github.com/spf13/cobra"
	"go-rest/internal/runtime"
	"go-rest/pkg/config"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the API server",
	Run: func(cmd *cobra.Command, args []string) {

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

		app := runtime.NewServerRuntime(cfg, logger)
		app.Run()
	},
}
