package infrastructure

import (
	"log"

	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // ✅ Assign DB globally
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("✅ Connected to Supabase via GORM!")

	if DB == nil {
		log.Fatal("❌ Database connection is NIL after initialization!")
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	log.Println("✅ Database migration successful!")
}
