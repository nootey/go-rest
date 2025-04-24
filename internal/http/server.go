package http

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	ginprometheus "github.com/zsais/go-gin-prometheus"
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

	router := NewRouter(container)

	addr := container.Config.Host + ":" + container.Config.Port

	return &Server{
		Router: router,
		logger: logger,
		server: &http.Server{
			Addr: addr,
		},
	}
}

func (s *Server) Start() {
	s.logger.Info("Starting the HTTP server")

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

	var r *gin.Engine
	var domainProtocol string

	if container.Config.Release {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		domainProtocol = "https://"

	} else {
		r = gin.Default()
		domainProtocol = "http://"
	}

	mongoCheck := checks.NewMongoCheck(10, container.DB)
	healthcheck.New(r, config.DefaultConfig(), []checks.Check{mongoCheck})

	// Setup CORS
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{
		domainProtocol + container.Config.WebClientDomain,
		domainProtocol + container.Config.WebClientDomain + ":" + container.Config.WebClientPort,
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "go-rest"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	p := ginprometheus.NewPrometheus("go_rest")
	p.Use(r)

	r.Use(gin.Recovery())

	routeInitializer := NewRouteInitializerHTTP(r, container)
	routeInitializer.InitEndpoints()

	return r
}
