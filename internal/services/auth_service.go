package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

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
		log.Printf("[ERROR] Supabase SignUp failed: %v", err)
		return "", "", "", errors.New("registration failed: " + err.Error())
	}

	if user == nil {
		log.Printf("[ERROR] Unexpected response from Supabase: user is nil")
		return "", "", "", errors.New("unexpected response from Supabase: user is nil")
	}

	log.Println("[INFO] Registration successful. Waiting for email confirmation.")
	fmt.Println("Please confirm your email before proceeding.")

	return user.ID, user.Email, "", nil
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

func (s *AuthService) UpdateUser(userID uuid.UUID, updates map[string]interface{}) error {
	return s.Repo.UpdateUser(userID, updates)
}

func (s *AuthService) LogoutUser(token string) error {
	if token == "" {
		log.Println("LogoutUser error: empty token")
		return errors.New("empty token")
	}

	err := s.SupabaseClient.Auth.SignOut(context.Background(), token)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			log.Println("Logout warning: token is already expired, treating as logged out.")
			return nil
		}
		log.Println("Supabase logout error:", err)
		return errors.New("failed to log out: " + err.Error())
	}

	return nil
}

func (s *AuthService) DeletePipeline(pipelineID string) error {
	if pipelineID == "" {
		return errors.New("pipeline ID cannot be empty")
	}

	err := s.Repo.DeletePipeline(context.Background(), pipelineID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete pipeline %s: %v", pipelineID, err)
		return err
	}

	log.Printf("[INFO] Pipeline %s deleted successfully", pipelineID)
	return nil
}
