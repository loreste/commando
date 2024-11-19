package models

import (
	"errors"
	"teltech/database"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the TelTech application
type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"size:100;unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"type:enum('user', 'admin');not null"`
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the stored password hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// FindByUsername finds a user by their username
func FindByUsername(username string) (*User, error) {
	var user User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func CreateUser(username, password, role string) (*User, error) {
	// Check if the user already exists
	existingUser, _ := FindByUsername(username)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Create new user
	user := &User{
		Username: username,
		Password: password,
		Role:     role,
	}

	// Hash the password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
