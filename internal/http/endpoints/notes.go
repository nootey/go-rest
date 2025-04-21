package endpoints

import (
	"github.com/gin-gonic/gin"
	httpHandlers "go-rest/internal/handlers"
)

func NotesRoutes(apiGroup *gin.RouterGroup, handler *httpHandlers.NotesHandler) {
	apiGroup.GET("/get", handler.GetNotes)
	apiGroup.POST("/create", handler.CreateNote)

}
