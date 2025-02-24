package services

import (
	"errors"

	"github.com/sarikap9/my-pipeline-project/internal/infrastructure"
	"github.com/sarikap9/my-pipeline-project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser creates a new user
func RegisterUser(email, password string) error {
	var existingUser models.User
	err := infrastructure.DB.QueryRow("SELECT id FROM users WHERE email=$1", email).Scan(&existingUser.ID)
	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = infrastructure.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, hashedPassword)
	return err
}

// AuthenticateUser verifies user credentials
func AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User
	err := infrastructure.DB.QueryRow("SELECT id, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	user.Email = email
	return &user, nil
}
