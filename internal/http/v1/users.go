package v1

import (
	"github.com/gin-gonic/gin"
	httpHandlers "go-rest/internal/http/handlers"
)

func UserRoutes(apiGroup *gin.RouterGroup, handler *httpHandlers.UserHandler) {
	apiGroup.GET("/get", handler.GetUsers)
	apiGroup.POST("/create", handler.CreateUser)

}
