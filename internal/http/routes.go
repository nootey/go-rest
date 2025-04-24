package http

import (
	"github.com/gin-gonic/gin"
	"go-rest/internal/bootstrap"
	httpHandlers "go-rest/internal/http/handlers"
	v1 "go-rest/internal/http/v1"
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
	api := r.Router.Group("/api")

	// Version 1
	_v1 := api.Group("/v1")
	r.initV1Routes(_v1)
}

func (r *RouteInitializerHTTP) initV1Routes(_v1 *gin.RouterGroup) {

	r.Router.GET("/", rootHandler)

	noteHandler := httpHandlers.NewNoteHandler(r.Container.NotesService)
	userHandler := httpHandlers.NewUserHandler(r.Container.UserService)
	authHandler := httpHandlers.NewAuthHandler(r.Container.AuthService)

	authGroup := _v1.Group("/", r.Container.AuthService.WebClientMiddleware.WebClientAuthentication())
	{
		userRoutes := authGroup.Group("/users")
		v1.UserRoutes(userRoutes, userHandler)
	}

	// Public routes
	publicGroup := _v1.Group("")
	{
		publicAuthRoutes := publicGroup.Group("/auth")
		v1.PublicAuthRoutes(publicAuthRoutes, authHandler)

		// These are public, just as an example
		notesRoutes := publicGroup.Group("/notes")
		v1.NotesRoutes(notesRoutes, noteHandler)

	}

}

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "HTTP server is running!"})
}
