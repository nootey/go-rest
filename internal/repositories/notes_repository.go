package repositories

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go-rest/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotesRepository struct {
	DB *mongo.Client
}

func NewNotesRepository(db *mongo.Client) *NotesRepository {
	return &NotesRepository{DB: db}
}

func (repo *NotesRepository) CreateNewNote(note *models.Note) error {

	if note == nil {
		return fmt.Errorf("note is nil")
	}

	existingNote := &models.Note{}
	err := mgm.Coll(note).First(bson.M{"title": note.Title}, existingNote)

	if err == nil {
		return fmt.Errorf("note with title %s already exists", note.Title)
	}

	return mgm.Coll(note).Create(note)
}

func (repo *NotesRepository) GetAllNotes() ([]models.Note, error) {
	var notes []models.Note
	err := mgm.Coll(&models.Note{}).SimpleFind(&notes, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving notes: %v", err)
	}

	return notes, err
}
