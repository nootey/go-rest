package http

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go-rest/internal/bootstrap"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	Router *gin.Engine
	server *http.Server
	logger *zap.Logger
}

func NewServer(container *bootstrap.Container, logger *zap.Logger) *Server {
	// Create a Router and attach middleware
	router := NewRouter(container)

	return &Server{
		Router: router,
		logger: logger,
		server: &http.Server{
			Addr: ":" + container.Config.Port,
		},
	}
}

func (s *Server) Start() {
	s.logger.Info("Starting the HTTP server")

	// Attach recovery & log middleware
	s.Router.Use(ginzap.Ginzap(s.logger, time.RFC3339, true), ginzap.RecoveryWithZap(s.logger, true))

	s.server.Handler = s.Router.Handler()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("Failed to listen and serve for: HTTP", zap.Error(err))
		}
	}()
}

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down the HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func NewRouter(container *bootstrap.Container) *gin.Engine {

	var router *gin.Engine
	var domainProtocol string

	if container.Config.Release {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		domainProtocol = "https://"

	} else {
		router = gin.Default()
		domainProtocol = "http://"
	}

	// Setup CORS
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{
		domainProtocol + container.Config.WebClientDomain,
		domainProtocol + container.Config.WebClientDomain + ":" + container.Config.WebClientPort,
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "go-rest"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Global middlewares
	router.Use(gin.Recovery())

	// Create RouteInitializer and initialize v1
	routeInitializer := NewRouteInitializerHTTP(router, container)
	routeInitializer.InitEndpoints()

	return router
}
