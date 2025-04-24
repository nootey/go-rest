package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
)

var rootCmd = &cobra.Command{
	Use:     "go-rest",
	Short:   "go-rest API server",
	Version: "1.0.0",
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Fatal("Failed to execute root command", zap.Error(err))
		log.Fatal(err)
	}
}
