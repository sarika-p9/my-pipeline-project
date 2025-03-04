package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
	pipeline_proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/pipeline"
	"github.com/sarika-p9/my-pipeline-project/api/http/handlers"
	"github.com/sarika-p9/my-pipeline-project/internal/adapters/primary"
	"github.com/sarika-p9/my-pipeline-project/internal/adapters/secondary"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/messaging"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
)

var (
	rabbitMQConn    *amqp.Connection
	rabbitMQChannel *amqp.Channel
	natsConn        *nats.Conn
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Proceeding with existing environment variables.")
	}

	// Initialize database and Supabase client
	infrastructure.InitDatabase()
	dbRepo := secondary.NewDatabaseAdapter()
	authService := services.NewAuthService(dbRepo)
	pipelineService := services.NewPipelineService(dbRepo)

	// Initialize REST API handler
	handler := &handlers.PipelineHandler{Service: pipelineService}
	authHandler := &handlers.AuthHandler{Service: authService}

	// Setup Gin router
	r := gin.Default()

	// Enable CORS for frontend integration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allows requests from any origin (including Postman Web & React)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/register", gin.WrapF(authHandler.RegisterHandler))
	r.POST("/login", gin.WrapF(authHandler.LoginHandler))
	r.POST("/pipelines", handler.CreatePipeline)
	r.POST("/pipelines/:id/start", handler.StartPipeline)
	r.GET("/pipelines/:id/status", handler.GetPipelineStatus)
	r.POST("/pipelines/:id/cancel", handler.CancelPipeline)
	r.GET("/user", gin.WrapF(authHandler.GetUserHandler))

	// Create gRPC server
	grpcServer := grpc.NewServer()
	authServer := &primary.AuthServer{AuthService: authService}
	pipelineServer := &primary.PipelineServer{Service: pipelineService}

	// Register gRPC services
	proto.RegisterAuthServiceServer(grpcServer, authServer)
	pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
	reflection.Register(grpcServer)

	// Start gRPC server
	go func() {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on port 50051: %v", err)
		}
		log.Println("Starting gRPC server on port 50051...")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Initialize RabbitMQ
	go func() {
		rabbitURL := "amqp://guest:guest@localhost:5672/"
		var err error
		rabbitMQConn, rabbitMQChannel, err = messaging.ConnectRabbitMQ(rabbitURL)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
		defer rabbitMQChannel.Close()
		defer rabbitMQConn.Close()
	}()

	// Initialize NATS
	go func() {
		natsURL := "nats://localhost:4222"
		natsConn = infrastructure.ConnectNATS(natsURL)
		defer natsConn.Close()
	}()

	// Start REST API server
	go func() {
		log.Println("Starting API server on port 8080...")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatalf("Failed to start REST API server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcServer.GracefulStop()
	log.Println("Servers exited gracefully")
}
