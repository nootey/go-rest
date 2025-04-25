package handlers

import (
	"github.com/gin-gonic/gin"
	"go-rest/internal/models"
	"go-rest/internal/services"
	"go-rest/pkg/utils"
	"net/http"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		Service: userService,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		utils.ErrorMessage(c, "Fetch error", err.Error(), http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorMessage(c, "Fetch error", err.Error(), http.StatusBadRequest, err)
		return
	}

	user.ModelVersion = models.UserModelVersion
	err := h.Service.CreateUser(&user)
	if err != nil {
		utils.ErrorMessage(c, "Create error", err.Error(), http.StatusBadRequest, err)
		return
	}

	utils.SuccessMessage(c, user.Email, "Create success", http.StatusOK)
}
