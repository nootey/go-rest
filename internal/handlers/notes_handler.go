package handlers

import (
	"go-rest/internal/models"
	"go-rest/internal/services"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotesHandler struct {
	service *services.NotesService
}

func NewNoteHandler(service *services.NotesService) *NotesHandler {
	return &NotesHandler{service: service}
}

func (h *NotesHandler) CreateNote(c *gin.Context) {

	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request["title"] == nil || request["description"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "note data is incomplete"})
		return
	}

	note := models.NewNote(request["title"].(string), request["description"].(string))

	err := h.service.Repo.CreateNewNote(note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bson.M{"Status": "Note has been created successfully!"})
}

func (h *NotesHandler) GetNotes(c *gin.Context) {
	notes, err := h.service.Repo.GetAllNotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}
