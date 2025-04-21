package bootstrap

import (
	"go-rest/internal/repositories"
	"go-rest/internal/services"
	"go-rest/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	Config       *config.Config
	DB           *mongo.Client
	NotesService *services.NotesService
}

func NewContainer(cfg *config.Config, db *mongo.Client) *Container {

	notesRepo := repositories.NewNotesRepository()

	notesService := services.NewtNotesService(notesRepo)

	return &Container{
		Config:       cfg,
		DB:           db,
		NotesService: notesService,
	}
}
