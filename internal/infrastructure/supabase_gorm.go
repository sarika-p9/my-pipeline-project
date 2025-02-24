package infrastructure

import (
	"fmt"
	"log"

	"github.com/sarikap9/my-pipeline-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitSupabaseWithGORM connects to Supabase and performs auto-migration for the Pipeline model.
func InitSupabaseWithGORM(supabaseURL, supabaseKey string) (*gorm.DB, error) {
	// Construct the DSN with parameters to disable prepared statement caching.
	// Make sure supabaseURL is your base connection string, e.g.:
	// postgresql://postgres.iizmevhufqeohsqxlxcm:your_password@aws-0-ap-south-1.pooler.supabase.com:6543/postgres
	dsn := fmt.Sprintf("%s?sslmode=require&prepareThreshold=0&prefer_simple_protocol=true", supabaseURL)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to Supabase: %v", err)
	}

	log.Println("✅ Connected to Supabase via GORM!")

	// Use a raw SQL query to check if the 'pipelines' table exists.
	var count int64
	// Note: table_name comparisons in Postgres are case-sensitive.
	err = db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", "pipelines").Scan(&count).Error
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to check table existence: %v", err)
	}
	if count > 0 {
		log.Println("✅ Pipeline table already exists, skipping auto migration.")
	} else {
		if err := db.AutoMigrate(&models.Pipeline{}); err != nil {
			return nil, fmt.Errorf("❌ AutoMigrate failed: %v", err)
		}
		log.Println("✅ Pipeline table auto-migrated!")
	}

	return db, nil
}
