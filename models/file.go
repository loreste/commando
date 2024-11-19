package models

import "time"

// File represents a file in the system
type File struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`        // Name of the file
	Path      string    `gorm:"unique;not null"` // File path in the file system
	Size      int64     `gorm:"not null"`        // File size in bytes
	FolderID  int       `gorm:"not null"`        // Foreign key to the parent folder
	CreatedAt time.Time `gorm:"autoCreateTime"`  // Timestamp when the file was created
	UpdatedAt time.Time `gorm:"autoUpdateTime"`  // Timestamp when the file was last updated
}
