package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Check normal auth
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ValidateJWT(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// check for valid admin token
func JWTAuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ValidateJWT(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		error := ValidateAdminRoleJWT(c)
		if error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only Administrator is allowed to perform this action"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// check for valid customer token
func JWTAuthMod() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		error := ValidateModRoleJWT(context)
		if error != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Only MODs are allowed to perform this action"})
			context.Abort()
			return
		}
		context.Next()
	}
}
