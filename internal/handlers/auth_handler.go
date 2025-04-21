package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-rest/internal/models"
	"go-rest/internal/services"
	"go-rest/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: authService,
	}
}

func (h *AuthHandler) GetAuthUser(c *gin.Context) {
	user, err := h.Service.GetCurrentUser(c)
	if err != nil {
		utils.ErrorMessage("Error occurred", err.Error(), http.StatusInternalServerError)(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var loginForm models.LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		utils.ErrorMessage("Error occurred", err.Error(), http.StatusBadRequest)(c, err)
		return
	}

	user, err := h.Service.UserRepo.GetUserByEmail(loginForm.Email, true)
	if err != nil || user == nil {
		utils.ErrorMessage("Error occurred", "Incorrect credentials", http.StatusUnauthorized)(c, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		utils.ErrorMessage("Error occurred", "Incorrect credentials", http.StatusUnauthorized)(c, err)
		return
	}

	userIDString := user.ID.Hex()

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.TokenData{
		UserID: userIDString,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-rest",
		},
	}

	jwtKey := []byte(h.Service.Config.JwtWebClientAccess)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		utils.ErrorMessage("Error occurred", "JWT generation failed", http.StatusInternalServerError)(c, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user.Email,
		"token": signedToken,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	_, claims, err := h.Service.WebClientMiddleware.ParseJWTToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	jwtSecret := []byte(h.Service.Config.JwtWebClientAccess)

	expirationTime := time.Now().Add(5 * time.Minute)
	newClaims := &models.TokenData{
		UserID: claims.UserID,
		Email:  claims.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-rest",
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newTokenString})
}

func (h *AuthHandler) LogoutUser(c *gin.Context) {
	utils.SuccessMessage("", "Logged out", http.StatusOK)(c.Writer, c.Request)
}
