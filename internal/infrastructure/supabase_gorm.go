package infrastructure

import (
	"log"
	"os"

	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := os.Getenv("SUPABASE_URL") // Ensure the correct env variable is used
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Ensure the necessary tables exist
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database connected successfully")
}

func GetDB() *gorm.DB {
	return DB
}
