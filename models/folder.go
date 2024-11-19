package models

import (
	"errors"
	"teltech/database"
)

// Folder represents a folder in the TelTech application
type Folder struct {
	ID      int    `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"size:255;not null"`
	Path    string `gorm:"size:255;unique;not null"`
	OwnerID int    `gorm:"not null"`
}

// CreateFolder saves a new folder to the database
func CreateFolder(name, path string, ownerID int) (*Folder, error) {
	// Check if a folder with the same path already exists
	existingFolder := Folder{}
	if err := database.DB.Where("path = ?", path).First(&existingFolder).Error; err == nil {
		return nil, errors.New("folder already exists at this path")
	}

	// Create the new folder
	folder := Folder{
		Name:    name,
		Path:    path,
		OwnerID: ownerID,
	}

	if err := database.DB.Create(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}

// GetFolderByPath retrieves a folder by its path
func GetFolderByPath(path string) (*Folder, error) {
	var folder Folder
	if err := database.DB.Where("path = ?", path).First(&folder).Error; err != nil {
		return nil, errors.New("folder not found")
	}
	return &folder, nil
}

// GetFolderByID retrieves a folder by its ID
func GetFolderByID(id int) (*Folder, error) {
	var folder Folder
	if err := database.DB.First(&folder, id).Error; err != nil {
		return nil, errors.New("folder not found")
	}
	return &folder, nil
}

// RenameFolder renames a folder in the database
func RenameFolder(folderID int, newName string) error {
	var folder Folder
	if err := database.DB.First(&folder, folderID).Error; err != nil {
		return errors.New("folder not found")
	}

	newPath := folder.Path[:len(folder.Path)-len(folder.Name)] + newName

	// Update the folder name and path
	folder.Name = newName
	folder.Path = newPath

	return database.DB.Save(&folder).Error
}

// DeleteFolder deletes a folder from the database
func DeleteFolder(folderID int) error {
	return database.DB.Delete(&Folder{}, folderID).Error
}

// GetFoldersByOwner retrieves all folders owned by a specific user
func GetFoldersByOwner(ownerID int) ([]Folder, error) {
	var folders []Folder
	if err := database.DB.Where("owner_id = ?", ownerID).Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}
