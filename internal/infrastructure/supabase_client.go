package infrastructure

import (
	// "log"
	"os"

	// "github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

func InitSupabaseClient() *supabase.Client {
	// if os.Getenv("K8S_ENV") != "true" { // Only load .env if not in Kubernetes
	// 	if err := godotenv.Load(); err != nil {
	// 		log.Println("No .env file found, using environment variables")
	// 	}
	// }

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	return supabase.CreateClient(url, key)
}
