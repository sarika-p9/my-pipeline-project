package infrastructure

import (
	"fmt"
	"log"
	"os"

	// "github.com/joho/godotenv"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// if os.Getenv("K8S_ENV") != "true" { // Only load .env if not in Kubernetes
	// 	if err := godotenv.Load(); err != nil {
	// 		log.Println("No .env file found, using environment variables")
	// 	}
	// }

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable is not set")
	}
	fmt.Println("Using DSN:", dsn)
	log.Printf("Connecting to database: %s", dsn)
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Supabase database: %v", err)
	}
	log.Println("Database connection established.")
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}
	log.Println("UUID-OSSP extension enabled.")
	migrateDatabase()
}

func migrateDatabase() {
	log.Println("Starting database migration...")
	migrateTable(&models.User{})
	migrateTable(&models.Pipelines{})
	migrateTable(&models.Stages{})
	log.Println("Database migration completed successfully.")
}

func migrateTable(model interface{}) {
	tableName := DB.NamingStrategy.TableName(getTableName(model))

	if DB.Migrator().HasTable(model) {
		log.Printf("Skipping migration: Table %s already exists.", tableName)
		return
	}

	if err := DB.AutoMigrate(model); err != nil {
		log.Fatalf("Failed to migrate table %s: %v", tableName, err)
	}
	log.Printf("Successfully migrated table: %s", tableName)
}

func getTableName(model interface{}) string {
	if t, ok := model.(interface{ TableName() string }); ok {
		return t.TableName()
	}
	return fmt.Sprintf("%T", model)
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database is not initialized. Call InitDatabase() first.")
	}
	return DB
}
