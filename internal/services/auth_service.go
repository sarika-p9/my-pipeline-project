package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"github.com/sarika-p9/my-pipeline-project/internal/ports"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

// AuthService handles user authentication
type AuthService struct {
	SupabaseClient *supabase.Client
	Repo           ports.PipelineRepository
}

// NewAuthService initializes the AuthService
func NewAuthService(repo ports.PipelineRepository) *AuthService {
	return &AuthService{
		SupabaseClient: infrastructure.InitSupabaseClient(), // Use the initialized Supabase client
		Repo:           repo,
	}
}

// RegisterUser registers a new user and waits for email confirmation
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

// LoginUser authenticates a user, ensures they exist in DB, and returns ID, email, and token
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
		}
		if err := s.Repo.SaveUser(newUser); err != nil {
			return "", "", "", errors.New("failed to save user in the database: " + err.Error())
		}
	}

	return session.User.ID, session.User.Email, session.AccessToken, nil
}
