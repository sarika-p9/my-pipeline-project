package infrastructure

import (
	"log"
	"os"

	"github.com/sarika-p9/my-pipeline-project/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDatabase sets up the database connection using environment variables
func InitDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is not set in environment variables")
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected")

	// Auto-migrate all necessary tables
	err = db.AutoMigrate(
		&types.User{}, // Add other models if needed
	)
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
}

// GetDatabaseInstance returns the initialized database connection
func GetDatabaseInstance() *gorm.DB {
	if db == nil {
		log.Fatal("❌ Database instance is not initialized. Call InitDatabase first.")
	}
	return db
}

// InsertUserIntoDB inserts a user into the database
func InsertUserIntoDB(user types.User) error {
	db := GetDatabaseInstance()
	return db.Create(&user).Error
}
