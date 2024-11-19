package controllers

import (
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"teltech/database"
	"teltech/models"

	"github.com/gin-gonic/gin"
)

var parentFolder = os.Getenv("PARENT_FOLDER") // Load parent folder from .env

// RenderFolder lists the contents of a folder
func RenderFolder(c *gin.Context) {
	folderPath := c.Query("path") // Optional query param for folder navigation
	if folderPath == "" {
		folderPath = parentFolder
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read folder contents"})
		return
	}

	contents := []map[string]string{}
	for _, file := range files {
		contents = append(contents, map[string]string{
			"name": file.Name(),
			"type": func() string {
				if file.IsDir() {
					return "folder"
				}
				return "file"
			}(),
		})
	}

	c.HTML(http.StatusOK, "folder.html", gin.H{
		"Title":   "TelTech File Manager",
		"Folders": contents,
		"Path":    folderPath,
	})
}

// CreateFolder creates a new folder
func CreateFolder(c *gin.Context) {
	var input struct {
		FolderName string `json:"folder_name" binding:"required"`
		ParentPath string `json:"parent_path"`
	}

	userID := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ParentPath == "" {
		input.ParentPath = parentFolder
	}

	fullPath := filepath.Join(input.ParentPath, input.FolderName)

	// Check if folder already exists
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder already exists"})
		return
	}

	// Create the folder
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	// Set ownership
	currentUser, err := user.LookupId(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lookup user"})
		return
	}

	uid, _ := strconv.Atoi(currentUser.Uid)
	gid, _ := strconv.Atoi(currentUser.Gid)

	if err := SetOwnership(fullPath, uid, gid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set folder ownership"})
		return
	}

	// Save folder metadata
	newFolder := models.Folder{
		Name:    input.FolderName,
		Path:    fullPath,
		OwnerID: userID,
	}
	if err := database.DB.Create(&newFolder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save folder metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder created successfully", "path": fullPath})
}

// RenameFolder renames an existing folder
func RenameFolder(c *gin.Context) {
	var input struct {
		FolderID int    `json:"folder_id" binding:"required"`
		NewName  string `json:"new_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the folder from the database
	var folder models.Folder
	if err := database.DB.First(&folder, input.FolderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	// Update folder name and path
	oldPath := folder.Path
	newPath := filepath.Join(filepath.Dir(oldPath), input.NewName)

	if err := os.Rename(oldPath, newPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename folder"})
		return
	}

	// Update database entry
	folder.Name = input.NewName
	folder.Path = newPath
	if err := database.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder renamed successfully"})
}

// DeleteFolder deletes a folder and its contents
func DeleteFolder(c *gin.Context) {
	var input struct {
		FolderID int `json:"folder_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the folder from the database
	var folder models.Folder
	if err := database.DB.First(&folder, input.FolderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	// Delete the folder from the file system
	if err := os.RemoveAll(folder.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}

	// Remove the folder from the database
	if err := database.DB.Delete(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}

// UploadFile uploads a file to a folder
func UploadFile(c *gin.Context) {
	parentPath := c.PostForm("parent_path")
	if parentPath == "" {
		parentPath = parentFolder
	}

	userID := c.GetInt("user_id")

	// Check if parent folder exists
	var folder models.Folder
	if err := database.DB.Where("path = ?", parentPath).First(&folder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent folder does not exist"})
		return
	}

	if folder.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this folder"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	uploadPath := filepath.Join(parentPath, file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Set ownership
	currentUser, err := user.LookupId(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lookup user"})
		return
	}

	uid, _ := strconv.Atoi(currentUser.Uid)
	gid, _ := strconv.Atoi(currentUser.Gid)

	if err := SetOwnership(uploadPath, uid, gid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set file ownership"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "path": uploadPath})
}

// DownloadFile serves a file as a download
func DownloadFile(c *gin.Context) {
	filePath := c.Query("file_path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file as a download
	c.FileAttachment(filePath, filepath.Base(filePath))
}
