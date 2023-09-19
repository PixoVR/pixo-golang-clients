package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func APIKeyAuthMiddleware(getCurrentUser func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthenticated := false

		token := ExtractToken(c.Request)

		if IsValidSecretKey(token) {
			isAuthenticated = true

		} else {
			var err error
			if err = getCurrentUser(c); err == nil {
				isAuthenticated = true
			}
		}

		if !isAuthenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Next()
	}
}
