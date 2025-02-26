package tests

import (
	"os"
	"testing"

	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
)

func TestInsertPipeline(t *testing.T) {
	// Initialize Supabase using environment variables
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		t.Fatal("❌ SUPABASE_URL or SUPABASE_KEY environment variables are not set")
	}

	infrastructure.InitSupabase(supabaseURL, supabaseKey)

	// Test inserting a pipeline
	err := infrastructure.InsertPipeline("Test Pipeline")
	if err != nil {
		t.Fatalf("❌ Failed to insert pipeline: %v", err)
	}

	t.Log("✅ Successfully inserted pipeline")
}

func TestGetPipelines(t *testing.T) {
	// Initialize Supabase using environment variables
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		t.Fatal("❌ SUPABASE_URL or SUPABASE_KEY environment variables are not set")
	}

	infrastructure.InitSupabase(supabaseURL, supabaseKey)

	// Test fetching pipelines
	pipelines, err := infrastructure.GetPipelines()
	if err != nil {
		t.Fatalf("❌ Failed to fetch pipelines: %v", err)
	}

	if len(pipelines) == 0 {
		t.Log("⚠️ No pipelines found")
	} else {
		t.Logf("✅ Retrieved %d pipelines", len(pipelines))
	}
}
