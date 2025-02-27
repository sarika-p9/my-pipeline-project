package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
)

// Register handles user registration via Supabase and stores the user in the database.
func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		log.Println("❌ Invalid request data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Println("🔹 Registering user:", req.Email)

	// Register the user in Supabase and store them in the database
	userID, err := services.RegisterUser(req.Email, req.Password)
	if err != nil {
		log.Println("❌ Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	log.Println("✅ User registered successfully:", req.Email, "with Supabase ID:", userID)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered. Please verify your email."})
}
