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
	// Build DSN without prepared statements.
	dsn := fmt.Sprintf("%s?sslmode=require", supabaseURL)

	// Open DB with GORM and enforce simple protocol
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // <-- Force simple protocol here
	}), &gorm.Config{
		PrepareStmt: false, // Disable GORM's prepared statements
	})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to Supabase: %v", err)
	}

	log.Println("✅ Connected to Supabase via GORM!")

	// Check if the 'pipelines' table exists.
	var count int64
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
