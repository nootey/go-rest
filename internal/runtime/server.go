package runtime

import (
	"context"
	"go-rest/internal/bootstrap"
	"go-rest/internal/http"
	"go-rest/pkg/config"
	"go-rest/pkg/database"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

type ServerRuntime struct {
	Logger *zap.Logger
	Config *config.Config
}

func NewServerRuntime(cfg *config.Config, logger *zap.Logger) *ServerRuntime {
	return &ServerRuntime{
		Config: cfg,
		Logger: logger,
	}
}

func (rt *ServerRuntime) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	dbClient, err := database.ConnectToDatabase(rt.Config)
	if err != nil {
		rt.Logger.Fatal("Database connection failed", zap.Error(err))
	}
	defer database.DisconnectDatabase()
	rt.Logger.Info("Successfully connected to the database")

	container := bootstrap.NewContainer(rt.Config, dbClient)
	httpLogger := rt.Logger.Named("http").With(zap.String("component", "HTTP"))
	httpServer := http.NewServer(container, httpLogger)
	go httpServer.Start()

	<-ctx.Done()

	rt.Logger.Info("Interrupt signal received, shutting down HTTP server...")
	if err := httpServer.Shutdown(); err != nil {
		rt.Logger.Fatal("HTTP server shutdown failed", zap.Error(err))
	}
	rt.Logger.Info("HTTP server exiting")
}
