package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	Service *services.AuthService
}

// Credentials struct for user input
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Received Register Request: Email=%s, Password=%s", creds.Email, creds.Password)

	userID, email, token, err := h.Service.RegisterUser(creds.Email, creds.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User Registered Successfully: ID=%s, Email=%s", userID, email)

	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
		"email":   email,
		"token":   token,
	})
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	log.Println("LoginHandler called") // Debugging log

	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Login request received: Email=%s", creds.Email)

	userID, email, token, err := h.Service.LoginUser(creds.Email, creds.Password)
	if err != nil {
		log.Println("Login failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Println("Login successful")
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
		"token":   token,
	})
}

// LogoutHandler handles user logout by revoking the session token
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Logout request binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.Service.LogoutUser(req.Token)
	if err != nil {
		log.Println("Logout error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
