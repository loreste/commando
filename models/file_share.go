package models

import "time"

// FileShare represents a shareable link for a file
type FileShare struct {
	ID         int        `gorm:"primaryKey;autoIncrement"`
	FileID     int        `gorm:"not null"`                                  // Foreign key for the file
	ShareLink  string     `gorm:"unique;not null"`                           // Unique shareable link
	AccessType string     `gorm:"type:enum('read', 'write');default:'read'"` // Access type: "read" or "write"
	Expiration *time.Time `gorm:"default:null"`                              // Optional expiration date
	Password   string     `gorm:"default:null"`                              // Optional password for the link
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
}
