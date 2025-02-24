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
	// Construct the DSN. Typically, Supabase provides a connection string,
	// but here we assume supabaseURL is the connection string without parameters.
	dsn := fmt.Sprintf("%s?sslmode=require", supabaseURL)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to Supabase: %v", err)
	}

	log.Println("✅ Connected to Supabase via GORM!")

	// Auto-migrate your Pipeline model.
	if err := db.AutoMigrate(&models.Pipeline{}); err != nil {
		return nil, fmt.Errorf("❌ AutoMigrate failed: %v", err)
	}
	log.Println("✅ Pipeline table auto-migrated!")

	return db, nil
}
