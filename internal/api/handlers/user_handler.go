package handlers

import (
	"fmt"
	"go-rest/internal/models"
	"go-rest/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{repository: repo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request["name"] == nil || request["email"] == nil || request["password"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User data is incomplete"})
		return
	}

	var user models.User
	user.Name = request["name"].(string)
	user.Email = request["email"].(string)
	user.Password = request["password"].(string)

	err := h.repository.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.repository.GetAllUsers()
	if err != nil {
		fmt.Println("Failed to get users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
