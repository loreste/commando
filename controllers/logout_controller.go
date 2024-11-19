package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout handles user logout
func Logout(c *gin.Context) {
	// For JWT-based authentication, instruct the client to delete the token
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful. Please remove your token on the client side.",
	})
}
