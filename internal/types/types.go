package types

type AuthResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID    string `gorm:"primaryKey"`
	Email string `gorm:"unique"`
}
