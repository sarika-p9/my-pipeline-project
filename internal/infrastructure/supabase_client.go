package infrastructure

import (
	"fmt"
	"log"

	"github.com/nedpals/supabase-go"
)

var SupabaseClient *supabase.Client

// Initialize Supabase client
func InitSupabase(url, key string) {
	SupabaseClient = supabase.CreateClient(url, key) // Fixed constructor
	log.Println("✅ Supabase client initialized.")
}

// InsertPipeline inserts a pipeline into the "pipelines" table
func InsertPipeline(pipelineName string) error {
	type Pipeline struct {
		Name string `json:"name"`
	}

	newPipeline := Pipeline{Name: pipelineName}

	var result []Pipeline

	// Corrected Insert call to handle response and error properly
	err := SupabaseClient.DB.From("pipelines").Insert(newPipeline).Execute(&result)
	if err != nil {
		return fmt.Errorf("❌ Failed to insert pipeline: %v", err)
	}

	log.Println("✅ Pipeline inserted successfully:", result)
	return nil
}

// GetPipelines retrieves all pipelines from the "pipelines" table
func GetPipelines() ([]map[string]interface{}, error) {
	var pipelines []map[string]interface{}

	// Select all columns without extra flags
	err := SupabaseClient.DB.From("pipelines").Select("*").Execute(&pipelines)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch pipelines: %v", err)
	}

	log.Printf("✅ Retrieved %d pipelines\n", len(pipelines))
	return pipelines, nil
}
