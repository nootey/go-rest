package http

import (
	"github.com/gin-gonic/gin"
	"go-rest/internal/bootstrap"
	"go-rest/internal/handlers"
	"go-rest/internal/http/endpoints"
	"go-rest/pkg/database"
	"net/http"
)

type RouteInitializerHTTP struct {
	Router    *gin.Engine
	Container *bootstrap.Container
}

func NewRouteInitializerHTTP(router *gin.Engine, container *bootstrap.Container) *RouteInitializerHTTP {
	return &RouteInitializerHTTP{
		Router:    router,
		Container: container,
	}
}

func (r *RouteInitializerHTTP) InitEndpoints() {
	apiPrefixV1 := "/api/v1"

	r.Router.GET("/", rootHandler)
	r.Router.GET(apiPrefixV1+"/health", func(c *gin.Context) {
		healthCheck(c)
	})

	noteHandler := handlers.NewNoteHandler(r.Container.NotesService)
	userHandler := handlers.NewUserHandler(r.Container.UserService)
	authHandler := handlers.NewAuthHandler(r.Container.AuthService)

	authGroup := r.Router.Group(apiPrefixV1, r.Container.AuthService.WebClientMiddleware.WebClientAuthentication())
	{
		userRoutes := authGroup.Group("/users")
		endpoints.UserRoutes(userRoutes, userHandler)
	}

	// Public routes
	publicGroup := r.Router.Group(apiPrefixV1)
	{
		publicAuthRoutes := publicGroup.Group("/auth")
		endpoints.PublicAuthRoutes(publicAuthRoutes, authHandler)

		// These are public, just as an example
		notesRoutes := publicGroup.Group("/notes")
		endpoints.NotesRoutes(notesRoutes, noteHandler)

	}

}

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "HTTP server is running!"})
}

func healthCheck(c *gin.Context) {
	httpHealthStatus := "healthy"
	dbStatus := "healthy"

	// Check database connection
	err := database.PingDatabase()
	if err != nil {
		dbStatus = "unhealthy"
		httpHealthStatus = "degraded"
	}

	statusCode := http.StatusOK
	if httpHealthStatus == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status": gin.H{
			"api": gin.H{"http": httpHealthStatus},
			"services": gin.H{
				"database": gin.H{"mongo": dbStatus},
			},
		},
	})
}
