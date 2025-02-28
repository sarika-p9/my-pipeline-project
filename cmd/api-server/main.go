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
	"github.com/streadway/amqp"
	"github.com/supabase-community/gotrue-go"
	"google.golang.org/grpc"

	"github.com/sarika-p9/my-pipeline-project/api/grpc/proto"
	"github.com/sarika-p9/my-pipeline-project/internal/adapter/secondary"
	"github.com/sarika-p9/my-pipeline-project/internal/grpcserver"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/messaging"
	"github.com/sarika-p9/my-pipeline-project/internal/middleware"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
)

var (
	grpcConn        *grpc.ClientConn
	grpcClient      proto.PipelineServiceClient
	rabbitMQConn    *amqp.Connection
	rabbitMQChannel *amqp.Channel
	natsConn        *nats.Conn
	supabaseAuth    *gotrue.Client
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

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Proceeding with existing environment variables.")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	postgresDSN := os.Getenv("POSTGRES_DSN") // Ensure this is set in the environment

	if supabaseURL == "" || supabaseKey == "" || postgresDSN == "" {
		log.Fatal("SUPABASE_URL, SUPABASE_KEY, and POSTGRES_DSN must be set in the environment variables")
	}

	// ✅ Initialize Supabase Client for Authentication
	infrastructure.InitSupabaseAuth(supabaseURL, supabaseKey)
	supabaseAuth = infrastructure.GetSupabaseAuth()

	// ✅ Initialize Database
	infrastructure.InitDatabase(postgresDSN)
	db := infrastructure.GetDB()

	initGRPC()
	defer closeGRPC()

	// ✅ Initialize database repository
	dbRepo := secondary.NewDatabaseAdapter(db)

	// ✅ Initialize authentication service
	authService := services.NewAuthService(dbRepo)

	// ✅ Initialize pipeline service
	pipelineService := services.NewPipelineService(dbRepo)

	// ✅ Initialize RabbitMQ
	go func() {
		rabbitURL := "amqp://guest:guest@localhost:5672/"
		var err error
		rabbitMQConn, rabbitMQChannel, err = messaging.ConnectRabbitMQ(rabbitURL)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
		defer func() {
			if rabbitMQChannel != nil {
				rabbitMQChannel.Close()
			}
			if rabbitMQConn != nil {
				rabbitMQConn.Close()
			}
		}()
	}()

	// ✅ Initialize NATS
	go func() {
		natsURL := "nats://localhost:4222"
		natsConn = infrastructure.ConnectNATS(natsURL)
		defer func() {
			if natsConn != nil {
				natsConn.Close()
			}
		}()
	}()

	services.SubscribeToEvents(natsConn, "pipeline.updates")

	go grpcserver.StartGRPCServer()

	r := gin.Default()

	// ✅ Initialize REST API handlers
	authHandler := services.NewAuthHandler(authService)
	pipelineHandler := services.NewPipelineHandler(pipelineService)

	// ✅ Public routes
	r.POST("/register", authHandler.RegisterHandler)
	r.POST("/login", authHandler.LoginHandler)

	// ✅ Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	authorized.POST("/pipelines", pipelineHandler.CreatePipeline)
	authorized.GET("/pipelines", pipelineHandler.ListPipelines)

	// ✅ API status and worker management
	r.GET("/status", services.GetAPIStatus)
	r.GET("/workers", services.ListWorkers)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// ✅ Start server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not run API server: %v", err)
		}
	}()

	// ✅ Graceful shutdown
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
