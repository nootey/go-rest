package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go-rest/internal/config"
	"go-rest/internal/handlers"
	"go-rest/internal/middleware"
	"go-rest/internal/repositories"
	"go-rest/internal/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
)

func main() {
	conf := config.LoadConfig()

	err := mgm.SetDefaultConfig(nil, "app", options.Client().ApplyURI(conf.MongoURI))
	if err != nil {
		log.Fatalf("Error initializing MongoDB: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logger.Sync()

	r := gin.Default()

	r.Use(middleware.Logger())

	userRepo := repositories.NewUserRepository()
	userHandler := handlers.NewUserHandler(userRepo)

	r.POST("/users/create", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.GET("/health", utils.HealthCheck)

	r.Run(":3000")
}
