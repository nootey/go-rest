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
		utils.ErrorMessage("fetch error", err.Error(), http.StatusBadRequest)(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorMessage("fetch error", err.Error(), http.StatusBadRequest)(c, err)
		return
	}

	err := h.Service.CreateUser(&user)
	if err != nil {
		utils.ErrorMessage("create error", err.Error(), http.StatusInternalServerError)(c, err)
		return
	}

	utils.SuccessMessage(user.Email, "User created successfully", http.StatusOK)(c.Writer, c.Request)
}
