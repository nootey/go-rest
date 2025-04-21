package repositories

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go-rest/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DB *mongo.Client
}

func NewUserRepository(db *mongo.Client) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User

	err := mgm.Coll(&models.User{}).SimpleFind(&users, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %v", err)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id string, includeSecrets bool) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	user := &models.User{}
	err = mgm.Coll(user).FindByID(objID, user)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) GetPasswordByEmail(email string) (string, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": email}, user)
	if err != nil {
		return "", fmt.Errorf("could not find user: %v", err)
	}

	return user.Password, nil
}

func (r *UserRepository) GetUserByEmail(email string, includeSecrets bool) (*models.User, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": email}, user)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	if !includeSecrets {
		user.Password = ""
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return mgm.Coll(user).Create(user)
}
