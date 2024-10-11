package repositories

import (
	"fmt"
	"go-rest/internal/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/kamva/mgm/v3"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) CreateUser(user *models.User) error {

	if user == nil {
		return fmt.Errorf("user is nil")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return mgm.Coll(user).Create(user)
}

func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := mgm.Coll(&models.User{}).SimpleFind(&users, nil)
	return users, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
