package services

import (
	"go-rest/internal/repositories"
)

type NotesService struct {
	Repo *repositories.NotesRepository
}

func NewtNotesService(
	repo *repositories.NotesRepository,
) *NotesService {
	return &NotesService{
		Repo: repo,
	}
}
