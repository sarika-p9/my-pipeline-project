package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
	"github.com/sarika-p9/my-pipeline-project/internal/core/ports"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

// AuthService handles user authentication
type AuthService struct {
	SupabaseClient *supabase.Client
	Repo           ports.PipelineRepository
}

// NewAuthService initializes the AuthService
func NewAuthService(repo ports.PipelineRepository) *AuthService {
	return &AuthService{
		SupabaseClient: infrastructure.InitSupabaseClient(),
		Repo:           repo,
	}
}

func (s *AuthService) RegisterUser(email, password string) (string, string, string, error) {
	user, err := s.SupabaseClient.Auth.SignUp(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", "", errors.New("registration failed: " + err.Error())
	}

	if user == nil {
		return "", "", "", errors.New("unexpected response from Supabase: user is nil")
	}

	// ✅ Wait for the email confirmation before proceeding
	fmt.Println("Please confirm your email before proceeding.")

	return user.ID, user.Email, "", nil // Don't return token yet
}

func (s *AuthService) LoginUser(email, password string) (string, string, string, error) {
	session, err := s.SupabaseClient.Auth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", "", errors.New("login failed: " + err.Error())
	}

	if session.User.ID == "" || session.User.Email == "" || session.AccessToken == "" {
		return "", "", "", errors.New("unexpected response from Supabase: session or user is nil")
	}

	// ✅ Convert string UserID to uuid.UUID
	userUUID, err := uuid.Parse(session.User.ID)
	if err != nil {
		return "", "", "", errors.New("invalid user UUID: " + err.Error())
	}

	// ✅ Check if user already exists in DB
	existingUser, _ := s.Repo.GetUserByID(userUUID)
	if existingUser == nil {
		// ✅ Save user details after email confirmation
		newUser := &models.User{
			UserID: userUUID,
			Email:  session.User.Email,
			// Name:   "Harsh Srivastava",
			// Role:   "admin",
		}
		if err := s.Repo.SaveUser(newUser); err != nil {
			return "", "", "", errors.New("failed to save user in the database: " + err.Error())
		}
	}

	return session.User.ID, session.User.Email, session.AccessToken, nil
}
