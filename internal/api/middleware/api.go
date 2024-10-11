package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func InitRateLimitMiddleware() gin.HandlerFunc {
	globalLimiter := rate.NewLimiter(rate.Every(1*time.Minute), 200)

	return func(c *gin.Context) {

		if !globalLimiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"title": "Too many requests", "message": "Try again later"})
			return
		}
		c.Next()

	}
}

func RateLimitMiddleware(requests int, timeframe time.Duration, message string) gin.HandlerFunc {
	if message == "" {
		message = "Try again later."
	}
	limiter := rate.NewLimiter(rate.Every(timeframe), requests)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"title": "Too many requests", "message": message})
			return
		}
		c.Next()
	}
}
