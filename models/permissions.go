package models

import (
	"errors"
	"teltech/database"
)

// Permission represents a user's access rights to a folder
type Permission struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	FolderID   int    `gorm:"not null"`
	UserID     int    `gorm:"not null"`
	Permission string `gorm:"type:enum('read', 'write', 'admin');not null"`
}

// AddPermission grants a user access to a folder with a specific permission level
func AddPermission(folderID, userID int, permission string) error {
	if permission != "read" && permission != "write" && permission != "admin" {
		return errors.New("invalid permission type")
	}

	existingPermission := Permission{}
	err := database.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).First(&existingPermission).Error
	if err == nil {
		return errors.New("user already has a permission for this folder")
	}

	newPermission := Permission{
		FolderID:   folderID,
		UserID:     userID,
		Permission: permission,
	}

	return database.DB.Create(&newPermission).Error
}

// UpdatePermission updates the permission level for a user on a folder
func UpdatePermission(folderID, userID int, newPermission string) error {
	if newPermission != "read" && newPermission != "write" && newPermission != "admin" {
		return errors.New("invalid permission type")
	}

	permission := Permission{}
	if err := database.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).First(&permission).Error; err != nil {
		return errors.New("permission not found for the given user and folder")
	}

	permission.Permission = newPermission
	return database.DB.Save(&permission).Error
}

// RemovePermission removes a user's access to a folder
func RemovePermission(folderID, userID int) error {
	return database.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).Delete(&Permission{}).Error
}

// GetUserPermission retrieves the permission level of a user for a specific folder
func GetUserPermission(folderID, userID int) (string, error) {
	permission := Permission{}
	if err := database.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).First(&permission).Error; err != nil {
		return "", errors.New("no permission found for this user on this folder")
	}
	return permission.Permission, nil
}

// HasPermission checks if a user has the required permission level for a folder
func HasPermission(folderID, userID int, requiredPermission string) (bool, error) {
	userPermission, err := GetUserPermission(folderID, userID)
	if err != nil {
		return false, err
	}

	// Permission hierarchy: admin > write > read
	permissionLevels := map[string]int{
		"read":  1,
		"write": 2,
		"admin": 3,
	}

	return permissionLevels[userPermission] >= permissionLevels[requiredPermission], nil
}
