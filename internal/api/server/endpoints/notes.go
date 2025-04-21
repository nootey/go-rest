package endpoints

import "github.com/gin-gonic/gin"
import httpHandlers "go-rest/internal/api/handlers"

func NotesRoutes(apiGroup *gin.RouterGroup, handler *httpHandlers.NotesHandler) {
	apiGroup.GET("/get", handler.GetNotes)
	apiGroup.POST("/create", handler.CreateNote)

}
