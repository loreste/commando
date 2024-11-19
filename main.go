package main

import (
	"log"
	"os"
	"teltech/models"

	"teltech/database"
	"teltech/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using defaults.")
	}

	// Initialize database
	database.InitDB()

	// Migrate database tables
	if err := database.DB.AutoMigrate(
		&models.User{},
		&models.Folder{},
		&models.File{},
		&models.FileShare{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Serve static files (CSS, JS, etc.)
	router.Static("/static", "./static")

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Setup application routes
	routes.SetupRoutes(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
