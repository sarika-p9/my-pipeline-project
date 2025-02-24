package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sarikap9/my-pipeline-project/api/grpc/proto"
	"github.com/sarikap9/my-pipeline-project/internal/grpcserver"
	"github.com/sarikap9/my-pipeline-project/internal/infrastructure" // Contains both raw Supabase and GORM code
	"google.golang.org/grpc"
)

var grpcConn *grpc.ClientConn
var grpcClient proto.PipelineServiceClient

// initGRPC initializes the gRPC client connection.
func initGRPC() {
	var err error
	grpcConn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	grpcClient = proto.NewPipelineServiceClient(grpcConn)
}

// closeGRPC closes the gRPC connection.
func closeGRPC() {
	if grpcConn != nil {
		grpcConn.Close()
	}
}

func main() {
	// Load environment variables from .env file.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve Supabase connection details from environment variables.
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set in the .env file")
	}

	// Initialize Supabase connection using GORM and perform auto-migration.
	db, err := infrastructure.InitSupabaseWithGORM(supabaseURL, supabaseKey)
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	log.Println("DB connection and migration successful:", db)

	// Start the gRPC server in a separate goroutine.
	go grpcserver.StartGRPCServer()

	// Initialize Gin router for the API server.
	r := gin.Default()

	// Initialize the gRPC client.
	initGRPC()
	defer closeGRPC()

	// API route: Check server status.
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API Server is running!"})
	})

	// API route: List pipelines via gRPC.
	r.GET("/pipelines", func(c *gin.Context) {
		resp, err := grpcClient.ListPipelines(context.Background(), &proto.ListPipelinesRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pipelines": resp.GetPipelines()})
	})

	// API route: Create new pipeline via gRPC.
	r.POST("/pipelines", func(c *gin.Context) {
		var newPipeline struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&newPipeline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := grpcClient.CreatePipeline(context.Background(), &proto.CreatePipelineRequest{
			Name: newPipeline.Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Pipeline created successfully!"})
	})

	// API route: Dummy worker list.
	r.GET("/workers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"workers": []string{"Worker1", "Worker2"}})
	})

	// Set up HTTP server.
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run the API server in a goroutine.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not run API server: %v", err)
		}
	}()

	// Graceful shutdown logic.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Create a timeout context for shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the API server gracefully.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("API server forced to shutdown: %v", err)
	}

	log.Println("Servers exited gracefully")
}
