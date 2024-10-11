package server

import (
	"context"
	"errors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go-rest/internal/api/middleware"
	goConfig "go-rest/pkg/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	Router *gin.Engine
	server *http.Server
	logger *zap.Logger
}

func NewServer(cfg *goConfig.Config, logger *zap.Logger) *Server {
	// Create a Router and attach middleware
	router := NewRouter(cfg)

	return &Server{
		Router: router,
		logger: logger.Named("http-server"),
		server: &http.Server{
			Addr: ":" + cfg.Port,
		},
	}
}

func (s *Server) Start() {
	s.logger.Info("Starting the server")

	// Attach recovery & log middleware
	s.Router.Use(ginzap.Ginzap(s.logger, time.RFC3339, true), ginzap.RecoveryWithZap(s.logger, true))

	s.server.Handler = s.Router.Handler()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("Failed to listen and serve", zap.Error(err))
		}
	}()
}

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func NewRouter(cfg *goConfig.Config) *gin.Engine {

	var router *gin.Engine

	if cfg.Release == "production" {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()

	} else {
		router = gin.Default()
	}

	// Global middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.InitRateLimitMiddleware())

	// Initialize API routes
	InitEndpoints(router, cfg)

	return router
}
