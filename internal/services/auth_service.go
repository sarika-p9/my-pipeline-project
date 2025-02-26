package services

import (
	"errors"
	"fmt"

	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser handles user registration
func RegisterUser(email, password string) error {
	var existingUser models.User

	// Check if user already exists
	err := infrastructure.DB.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create new user
	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := infrastructure.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// AuthenticateUser verifies user credentials
func AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User
	err := infrastructure.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err // Handle other DB errors
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	user.Email = email
	return &user, nil
}
