package handlers

import (
	"go-rest/internal/models"
	"go-rest/internal/services"
	"go-rest/pkg/utils"
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

	note := &models.Note{
		Title:       request["title"].(string),
		Description: request["description"].(string),
	}

	err := h.service.CreateNote(note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.SuccessMessage("Note has been created successfully!", "Create success", http.StatusOK)(c.Writer, c.Request)
}

func (h *NotesHandler) GetNotes(c *gin.Context) {
	notes, err := h.service.Repo.GetAllNotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}
