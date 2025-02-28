package infrastructure

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Using system environment variables.")
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable is not set")
	}

	log.Printf("Connecting to database: %s", dsn)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Supabase database: %v", err)
	}
	log.Println("Database connection established.")

	// Enable UUID extension
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}
	log.Println("UUID-OSSP extension enabled.")

	// **Migrate User table first**
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate User table: %v", err)
	}

	// **Migrate PipelineExecution next**
	if err := DB.AutoMigrate(&models.PipelineExecution{}); err != nil {
		log.Fatalf("Failed to migrate PipelineExecution table: %v", err)
	}

	// **Migrate ExecutionLog last (depends on PipelineExecution)**
	if err := DB.AutoMigrate(&models.ExecutionLog{}); err != nil {
		log.Fatalf("Failed to migrate ExecutionLog table: %v", err)
	}

	log.Println("Database migration completed successfully.")
}

// ✅ Add this function to return the initialized DB instance
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database is not initialized. Call InitDatabase() first.")
	}
	return DB
}
