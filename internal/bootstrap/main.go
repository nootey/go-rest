package bootstrap

import (
	"go-rest/internal/middleware"
	"go-rest/internal/repositories"
	"go-rest/internal/services"
	"go-rest/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	Config       *config.Config
	DB           *mongo.Client
	Middleware   *middleware.WebClientMiddleware
	NotesService *services.NotesService
	UserService  *services.UserService
	AuthService  *services.AuthService
}

func NewContainer(cfg *config.Config, db *mongo.Client) *Container {

	webClientMiddleware := middleware.NewWebClientMiddleware(cfg)

	notesRepo := repositories.NewNotesRepository(db)
	userRepo := repositories.NewUserRepository(db)

	notesService := services.NewNotesService(notesRepo)
	userService := services.NewUserService(cfg, userRepo)
	authService := services.NewAuthService(cfg, userRepo, webClientMiddleware)

	return &Container{
		Config:       cfg,
		DB:           db,
		Middleware:   webClientMiddleware,
		NotesService: notesService,
		UserService:  userService,
		AuthService:  authService,
	}
}
