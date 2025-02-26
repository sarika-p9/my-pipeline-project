package models

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}

// CreateUser inserts a new user using GORM.
func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

// GetUserByEmail retrieves a user by email using GORM.
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("No user found with email: %s", email)
		return nil, nil
	} else if err != nil {
		log.Printf("Error during query: %v", err)
		return nil, err
	}

	log.Printf("Found user: %+v", user)
	return &user, nil
}
