package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

func RegisterUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Ensure using global DB
	db := infrastructure.DB
	if db == nil {
		log.Println("❌ Database connection is NIL in RegisterUser!")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	// Step 3: Log DB connection
	log.Printf("Using DB connection: %v", db)

	// Step 4: Log existing users
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
	} else {
		log.Printf("Users in DB: %+v", users)
	}

	// Check if user exists
	existingUser, err := models.GetUserByEmail(db, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for existing user"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Create user
	if err := input.CreateUser(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
