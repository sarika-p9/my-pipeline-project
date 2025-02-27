package infrastructure

import (
	"log"

	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/supabase-go"
)

var (
	SupabaseAuth   gotrue.Client
	SupabaseClient *supabase.Client
)

func InitSupabase(url, key string) {
	// Initialize authentication client
	SupabaseAuth = gotrue.New(url, key)
	if SupabaseAuth == nil {
		log.Fatal("Failed to initialize Supabase authentication client")
	}

	// Initialize Supabase client (handling the error)
	var err error
	SupabaseClient, err = supabase.NewClient(url, key, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}
}

// Get the Supabase client
func GetSupabaseClient() *supabase.Client {
	if SupabaseClient == nil {
		log.Fatal("Supabase client is not initialized. Call InitSupabase first.")
	}
	return SupabaseClient
}
