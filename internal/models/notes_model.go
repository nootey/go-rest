package models

import (
	"github.com/kamva/mgm/v3"
)

type Note struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string `json:"title" bson:"title"`
	Description      string `json:"description" bson:"description"`
}

func NewNote(title string, description string) *Note {
	return &Note{
		Title:       title,
		Description: description,
	}
}
