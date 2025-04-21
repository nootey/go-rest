package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`

	FirstName     string     `bson:"first_name" json:"first_name" binding:"required"`
	LastName      string     `bson:"last_name" json:"last_name" binding:"required"`
	Password      string     `bson:"password" json:"-"`
	Email         string     `bson:"email" json:"email" binding:"required,email"`
	EmailVerified *time.Time `bson:"email_verified,omitempty" json:"email_verified,omitempty"`
	Role          string     `bson:"role" json:"role" binding:"required"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty" json:"-"`
}
