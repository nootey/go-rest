package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-rest/internal/models"
	"go-rest/pkg/config"
	"net/http"
	"strings"
)

var (
	ErrTokenExpired = errors.New("token has expired")
)

type WebClientUserClaim struct {
	UserID string `json:"ID"`
	jwt.RegisteredClaims
}

type WebClientMiddleware struct {
	config *config.Config
}

func NewWebClientMiddleware(cfg *config.Config) *WebClientMiddleware {
	return &WebClientMiddleware{
		config: cfg,
	}
}

func (m *WebClientMiddleware) WebClientAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		jwtSecret := []byte(m.config.JwtWebClientAccess)
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := &models.TokenData{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Expired token"})
			return
		}

		// Set user email in context
		c.Set("email", claims.Email)

		c.Next()
	}
}

func (m *WebClientMiddleware) ParseJWTToken(c *gin.Context) (*jwt.Token, *models.TokenData, error) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, nil, fmt.Errorf("no Authorization header found")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	claims := &models.TokenData{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.JwtWebClientAccess), nil
	})

	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf("invalid or expired token")
	}

	return token, claims, nil
}
