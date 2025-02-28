package infrastructure

import (
	"fmt"

	"github.com/sarika-p9/my-pipeline-project/internal/models"
)

// StoreUserInDB stores a user in PostgreSQL
func StoreUserInDB(user *models.User) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Create(user).Error
}

// UpdateUserInDB updates user details
func UpdateUserInDB(user *models.User) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Save(user).Error
}
