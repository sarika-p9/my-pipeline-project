package infrastructure

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sarikap9/my-pipeline-project/internal/models"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	log.Println("Connected to Supabase (PostgreSQL)")
}

// CreateUser inserts a new user into the database.
func CreateUser(user *models.User) error {
	_, err := DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

// GetUserByEmail retrieves a user by email.
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
