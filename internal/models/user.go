package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"` // Supabase User ID
	Email     string    `gorm:"unique;not null" json:"email"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func VerifyUser(db *gorm.DB, token string) error {
	return db.Model(&User{}).Where("verification_token = ?", token).Update("verified", true).Error
}
