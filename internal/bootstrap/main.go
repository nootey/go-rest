package bootstrap

import (
	"go-rest/internal/api/services"
	"go-rest/internal/repositories"
	"go-rest/pkg/config"
)

type Container struct {
	Config       *config.Config
	NotesService *services.NotesService
}

func NewContainer(cfg *config.Config) *Container {

	notesRepo := repositories.NewNotesRepository()

	notesService := services.NewtNotesService(notesRepo)

	return &Container{
		Config:       cfg,
		NotesService: notesService,
	}
}
