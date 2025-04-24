package v1

import (
	"github.com/gin-gonic/gin"
	httpHandlers "go-rest/internal/http/handlers"
)

func PublicAuthRoutes(apiGroup *gin.RouterGroup, handler *httpHandlers.AuthHandler) {
	apiGroup.POST("/login", handler.LoginUser)
	apiGroup.GET("/refresh_token", handler.RefreshToken)
	apiGroup.POST("/logout", handler.LogoutUser)
}
