package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
	"github.com/sarika-p9/my-pipeline-project/internal/core/ports"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

type AuthService struct {
	SupabaseClient *supabase.Client
	Repo           ports.PipelineRepository
}

func NewAuthService(repo ports.PipelineRepository) *AuthService {
	return &AuthService{
		SupabaseClient: infrastructure.InitSupabaseClient(),
		Repo:           repo,
	}
}

func (s *AuthService) RegisterUser(email, password string) (string, string, string, error) {
	log.Printf("[DEBUG] RegisterUser called with Email: %s", email)

	user, err := s.SupabaseClient.Auth.SignUp(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Printf("[ERROR] Supabase SignUp failed: %v", err) // Log error
		return "", "", "", errors.New("registration failed: " + err.Error())
	}

	if user == nil {
		log.Printf("[ERROR] Unexpected response from Supabase: user is nil") // Log if user is nil
		return "", "", "", errors.New("unexpected response from Supabase: user is nil")
	}

	// ✅ Wait for the email confirmation before proceeding
	log.Println("[INFO] Registration successful. Waiting for email confirmation.")
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

	userUUID, err := uuid.Parse(session.User.ID)
	if err != nil {
		return "", "", "", errors.New("invalid user UUID: " + err.Error())
	}

	existingUser, _ := s.Repo.GetUserByID(userUUID)
	if existingUser == nil {
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

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.Repo.GetUserByID(userID)
}

func (s *AuthService) GetUserByToken(token string) (*models.User, error) {
	// Validate token with Supabase
	user, err := s.SupabaseClient.Auth.User(context.Background(), token)
	if err != nil {
		log.Printf("[ERROR] Invalid token: %v", err)
		return nil, errors.New("invalid token")
	}

	// Parse UUID
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, errors.New("invalid user UUID")
	}

	// Fetch user details from database
	return s.Repo.GetUserByID(userUUID)
}
