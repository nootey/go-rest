package services

import (
	"errors"
	"go-rest/internal/models"
	"go-rest/internal/repositories"
	"go-rest/pkg/config"
	"go-rest/pkg/utils"
)

type UserService struct {
	UserRepo *repositories.UserRepository
	Config   *config.Config
}

func NewUserService(cfg *config.Config, userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		Config:   cfg,
	}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) CreateUser(user *models.User) error {

	if user.Password == "" {
		return errors.New("password must be included")
	}

	hashedPassword, err := utils.HashAndSaltPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
