package repositories

import (
	"fmt"
	"go-rest/internal/models"

	"github.com/kamva/mgm/v3"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) Create(user *models.User) error {

	if user == nil {
		return fmt.Errorf("user is nil")
	}

	return mgm.Coll(user).Create(user)
}

func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := mgm.Coll(&models.User{}).SimpleFind(&users, nil)
	return users, err
}
