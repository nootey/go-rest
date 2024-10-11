package server

import (
	"github.com/gin-gonic/gin"
	"go-rest/internal/api/handlers"
	"go-rest/internal/repositories"
	"go-rest/pkg/config"
	"net/http"
)

func checkState(c *gin.Context) {
	response := gin.H{
		"status":  "healthy",
		"message": "API is running",
		"code":    200,
	}

	c.JSON(http.StatusOK, response)
}

func InitEndpoints(r *gin.Engine, cfg *config.Config) {

	apiV1 := "/api/v1"

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Go api")
	})

	r.GET(apiV1+"/health", func(c *gin.Context) {
		checkState(c)
	})

	userRepo := repositories.NewUserRepository()
	userHandler := handlers.NewUserHandler(userRepo)

	apiGroup := r.Group(apiV1)
	{
		apiGroup.POST("/users/create", userHandler.CreateUser)
		apiGroup.GET("/users", userHandler.GetUsers)
	}

}
