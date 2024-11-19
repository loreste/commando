package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"teltech/database"
	"teltech/models"

	"github.com/gin-gonic/gin"
)

// GenerateShareableLink generates a link for sharing a file
func GenerateShareableLink(c *gin.Context) {
	var input struct {
		FileID     int    `json:"file_id" binding:"required"` // ID of the file to share
		AccessType string `json:"access_type"`                // "read" or "write"
		Expiration string `json:"expiration"`                 // Expiration date (optional, RFC3339 format)
		Password   string `json:"password"`                   // Password for protection (optional)
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the file exists
	var file models.File
	if err := database.DB.First(&file, input.FileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Generate a random share link
	shareLink := generateRandomLink()

	// Parse expiration date (if provided)
	var expiration *time.Time
	if input.Expiration != "" {
		exp, err := time.Parse(time.RFC3339, input.Expiration)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration format"})
			return
		}
		expiration = &exp
	}

	// Create the share record
	share := models.FileShare{
		FileID:     input.FileID,
		ShareLink:  shareLink,
		AccessType: input.AccessType,
		Expiration: expiration,
		Password:   input.Password,
	}

	if err := database.DB.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"share_link": shareLink,
		"message":    "Share link generated successfully",
	})
}

// AccessSharedFile allows users to access a file via a shareable link
func AccessSharedFile(c *gin.Context) {
	shareLink := c.Param("share_link")
	var share models.FileShare

	// Find the share record
	if err := database.DB.Where("share_link = ?", shareLink).First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share link not found"})
		return
	}

	// Check expiration
	if share.Expiration != nil && time.Now().After(*share.Expiration) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Share link has expired"})
		return
	}

	// Check password (if set)
	if share.Password != "" {
		providedPassword := c.Query("password")
		if providedPassword != share.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}
	}

	// Find the file and serve it
	var file models.File
	if err := database.DB.First(&file, share.FileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file
	c.File(file.Path)
}

// generateRandomLink generates a secure random string for the share link
func generateRandomLink() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
