package routes

import (
	"teltech/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the routes for TelTech
func SetupRoutes(router *gin.Engine) {
	// Dashboard route
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", gin.H{
			"Title": "Dashboard",
		})
	})

	// File manager route
	router.GET("/file-manager", func(c *gin.Context) {
		c.HTML(200, "folder.html", gin.H{
			"Title": "File Manager",
		})
	})

	// Folder management routes
	router.POST("/folder/create", controllers.CreateFolder)   // Create a new folder
	router.PUT("/folder/rename", controllers.RenameFolder)    // Rename an existing folder
	router.DELETE("/folder/delete", controllers.DeleteFolder) // Delete a folder and its contents

	// File management routes
	router.POST("/file/upload", controllers.UploadFile)    // Upload a file
	router.GET("/file/download", controllers.DownloadFile) // Download a file

	// File sharing routes
	router.POST("/file/share", controllers.GenerateShareableLink)       // Generate a shareable link
	router.GET("/file/share/:share_link", controllers.AccessSharedFile) // Access a file via shareable link

	// Dashboard summary data
	router.GET("/api/dashboard/summary", controllers.GetDashboardSummary) // Get summary data for dashboard

	// Authentication routes
	router.POST("/register", controllers.Register) // Register a new user
	router.POST("/login", controllers.Login)       // Login for existing users
	router.POST("/logout", controllers.Logout)     // Logout the user
}
