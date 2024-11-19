package controllers

import (
	"net/http"
	"teltech/database"
	"teltech/models"

	"github.com/gin-gonic/gin"
)

// GetDashboardSummary returns a summary of dashboard statistics
func GetDashboardSummary(c *gin.Context) {
	var totalFolders int64
	var totalFiles int64
	var totalUsers int64

	// Count total folders
	if err := database.DB.Model(&models.Folder{}).Count(&totalFolders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch folder count"})
		return
	}

	// Count total files (assumes files are part of a separate model, e.g., File)
	if err := database.DB.Table("files").Count(&totalFiles).Error; err != nil {
		// Return 0 files if the table or data doesn't exist
		totalFiles = 0
	}

	// Count total users
	if err := database.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user count"})
		return
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"totalFolders": totalFolders,
		"totalFiles":   totalFiles,
		"totalUsers":   totalUsers,
	})
}
