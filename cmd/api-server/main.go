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
	"github.com/sarikap9/my-pipeline-project/api/grpc/proto" // Ensure this is the correct path
	"google.golang.org/grpc"
)

var grpcConn *grpc.ClientConn
var grpcClient proto.PipelineServiceClient

// Initialize gRPC client connection
func initGRPC() {
	var err error
	grpcConn, err = grpc.Dial("localhost:50051", grpc.WithInsecure()) // Connect to gRPC server
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	grpcClient = proto.NewPipelineServiceClient(grpcConn)
}

// Clean up gRPC connection
func closeGRPC() {
	if grpcConn != nil {
		grpcConn.Close()
	}
}

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize gRPC client
	initGRPC()
	defer closeGRPC()

	// Define routes
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API Server is running!"})
	})

	// Get list of pipelines from gRPC service
	r.GET("/pipelines", func(c *gin.Context) {
		resp, err := grpcClient.ListPipelines(context.Background(), &proto.ListPipelinesRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pipelines": resp.GetPipelines()})
	})

	// Create a new pipeline (Fixed field name)
	r.POST("/pipelines", func(c *gin.Context) {
		var newPipeline struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&newPipeline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := grpcClient.CreatePipeline(context.Background(), &proto.CreatePipelineRequest{
			Name: newPipeline.Name, // ✅ Correct field name
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Pipeline created successfully!"})
	})

	// Dummy worker list endpoint (Implement real logic later)
	r.GET("/workers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"workers": []string{"Worker1", "Worker2"}})
	})

	// Graceful shutdown handling
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	// Graceful shutdown logic
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
