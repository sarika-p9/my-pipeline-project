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
	"github.com/nats-io/nats.go"
	"github.com/sarika-p9/my-pipeline-project/api/grpc/proto"
	"github.com/sarika-p9/my-pipeline-project/api/http/handlers"
	"github.com/sarika-p9/my-pipeline-project/internal/grpcserver"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/messaging"
	"github.com/sarika-p9/my-pipeline-project/internal/middleware"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

var (
	grpcConn        *grpc.ClientConn
	grpcClient      proto.PipelineServiceClient
	rabbitMQConn    *amqp.Connection
	rabbitMQChannel *amqp.Channel
	natsConn        *nats.Conn
)

// Initialize gRPC connection
func initGRPC() {
	var err error
	grpcConn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	grpcClient = proto.NewPipelineServiceClient(grpcConn)
}

func closeGRPC() {
	if grpcConn != nil {
		grpcConn.Close()
	}
}

// Initialize RabbitMQ
func initRabbitMQ() {
	var err error
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	rabbitMQConn, rabbitMQChannel, err = messaging.ConnectRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
}

func closeRabbitMQ() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
	if rabbitMQConn != nil {
		rabbitMQConn.Close()
	}
}

// Initialize NATS
func initNATS() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	natsConn = infrastructure.ConnectNATS(natsURL)
}

func closeNATS() {
	if natsConn != nil {
		natsConn.Close()
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set in the .env file")
	}

	// Initialize database connection
	infrastructure.InitDatabase()
	db := infrastructure.GetDatabaseInstance()

	initGRPC()
	defer closeGRPC()

	initRabbitMQ()
	defer closeRabbitMQ()

	initNATS()
	defer closeNATS()

	services.SubscribeToEvents(natsConn, "pipeline.updates")

	go grpcserver.StartGRPCServer()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Initialize handlers with DB instance
	authHandler := handlers.AuthHandler{DB: db}
	pipelineHandler := handlers.NewPipelineHandler(db)

	// Public routes
	r.POST("/register", authHandler.SignupHandler)
	r.POST("/login", authHandler.LoginHandler)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	authorized.POST("/pipelines", pipelineHandler.CreatePipeline)
	authorized.GET("/pipelines", pipelineHandler.ListPipelines)

	// API status and worker management
	r.GET("/status", handlers.GetAPIStatus)
	r.GET("/workers", handlers.ListWorkers)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not run API server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("API server forced to shutdown: %v", err)
	}

	log.Println("Servers exited gracefully")
}
