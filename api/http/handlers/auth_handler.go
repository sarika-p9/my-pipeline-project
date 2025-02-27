package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	DB *gorm.DB
}

// SignupHandler registers a user in Supabase and stores info in PostgreSQL
func (h *AuthHandler) SignupHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Register user in Supabase
	user, err := services.RegisterUserInSupabase(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure Supabase ID exists
	if user.ID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Supabase user ID"})
		return
	}

	// Store user info in PostgreSQL
	newUser := models.User{
		ID:       user.ID, // ✅ Use Supabase user ID
		Email:    req.Email,
		Verified: false, // Default to false until verified
	}

	if err := models.CreateUser(h.DB, &newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user in database"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"user": gin.H{"id": newUser.ID, "email": newUser.Email}})
}

// LoginHandler logs in a user via Supabase
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Authenticate user via Supabase
	session, err := services.LoginUserInSupabase(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"session": session})
}
