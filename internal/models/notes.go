package models

import (
	"github.com/kamva/mgm/v3"
)

type Note struct {
	mgm.DefaultModel `bson:",inline"`
	ModelVersion     int `bson:"model_version" json:"model_version"`
	
	Title       string `bson:"title" json:"title" binding:"required"`
	Description string `bson:"description" json:"description"`
}
