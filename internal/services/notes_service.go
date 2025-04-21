package services

import (
	"go-rest/internal/models"
	"go-rest/internal/repositories"
)

type NotesService struct {
	Repo *repositories.NotesRepository
}

func NewNotesService(
	repo *repositories.NotesRepository,
) *NotesService {
	return &NotesService{
		Repo: repo,
	}
}

func (s *NotesService) CreateNote(record *models.Note) error {

	err := s.Repo.InsertNote(record)
	if err != nil {
		return err
	}
	return nil
}
