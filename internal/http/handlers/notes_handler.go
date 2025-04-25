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
		utils.ErrorMessage(c, "Bind error", err.Error(), http.StatusBadRequest, err)
		return
	}

	note := &models.Note{
		ModelVersion: models.NotesModelVersion,
		Title:        request["title"].(string),
		Description:  request["description"].(string),
	}

	err := h.service.CreateNote(note)
	if err != nil {
		utils.ErrorMessage(c, "Create error", err.Error(), http.StatusBadRequest, err)
		return
	}

	utils.SuccessMessage(c, "Note has been created successfully!", "Create success", http.StatusOK)
}

func (h *NotesHandler) GetNotes(c *gin.Context) {
	notes, err := h.service.Repo.GetAllNotes()
	if err != nil {
		utils.ErrorMessage(c, "General error", err.Error(), http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, notes)
}
