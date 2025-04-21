package services

import (
	"github.com/gin-gonic/gin"
	"go-rest/internal/middleware"
	"go-rest/internal/models"
	"go-rest/internal/repositories"
	"go-rest/pkg/config"
)

type AuthService struct {
	Config              *config.Config
	UserRepo            *repositories.UserRepository
	WebClientMiddleware *middleware.WebClientMiddleware
}

func NewAuthService(
	cfg *config.Config,
	userRepo *repositories.UserRepository,
	webClientMiddleware *middleware.WebClientMiddleware,
) *AuthService {
	return &AuthService{
		Config:              cfg,
		UserRepo:            userRepo,
		WebClientMiddleware: webClientMiddleware,
	}
}

func (s *AuthService) GetCurrentUser(c *gin.Context) (*models.User, error) {
	_, claims, err := s.WebClientMiddleware.ParseJWTToken(c)
	if err != nil {
		return nil, err
	}

	userID := claims.UserID

	user, err := s.UserRepo.GetUserByID(userID, false)
	if err != nil {
		return nil, err
	}

	return user, nil
}
