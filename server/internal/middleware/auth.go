// File: internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/JettZgg/LineUp/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		if uid, ok := claims["uid"].(float64); ok {
			c.Set("uid", int64(uid))
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user UID in token"})
			c.Abort()
			return
		}

		c.Next()
	}
}