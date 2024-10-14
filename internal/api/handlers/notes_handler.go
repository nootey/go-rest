package handlers

import (
	"go-rest/internal/models"
	"go-rest/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	repository *repositories.NoteRepository
}

func NewNoteHandler(repo *repositories.NoteRepository) *NoteHandler {
	return &NoteHandler{repository: repo}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {

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

	err := h.repository.CreateNewNote(note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bson.M{"Status": "Note has been created successfully!"})
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	notes, err := h.repository.GetAllNotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}
