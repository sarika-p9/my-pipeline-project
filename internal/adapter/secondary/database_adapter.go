package infrastructure

import (
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"gorm.io/gorm"
)

// Ensure User struct is from the models package
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
