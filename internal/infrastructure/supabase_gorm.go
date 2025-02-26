package infrastructure

import (
	"fmt"
	"log"
	"os"

	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB // Global DB variable

// InitSupabaseWithGORM connects to Supabase and performs auto-migration for the Pipeline model.
func InitSupabaseWithGORM(supabaseURL, supabaseKey string) (*gorm.DB, error) {
	// Build DSN without prepared statements.
	dsn := fmt.Sprintf("%s?sslmode=require", supabaseURL)

	// Set up GORM logger to log SQL queries.
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Log to stdout
		logger.Config{
			SlowThreshold:             200,         // Flag slow queries (>200ms)
			LogLevel:                  logger.Info, // Log all queries
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open DB with GORM and enforce simple protocol
	var err error
	dbInstance, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Force simple protocol here
	}), &gorm.Config{
		PrepareStmt: false,     // Disable GORM's prepared statements
		Logger:      newLogger, // Apply the logger here
	})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to Supabase: %v", err)
	}

	log.Println("✅ Connected to Supabase via GORM!")

	// Check if the 'pipelines' table exists.
	var count int64
	err = dbInstance.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", "pipelines").Scan(&count).Error
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to check table existence: %v", err)
	}
	if count > 0 {
		log.Println("✅ Pipeline table already exists, skipping auto migration.")
	} else {
		// Check and migrate the 'users' table
		err = dbInstance.AutoMigrate(&models.User{})
		if err != nil {
			log.Fatalf("AutoMigrate failed for users table: %v", err)
		}
		log.Println("✅ Users table auto-migrated!")
	}

	return dbInstance, nil
}
