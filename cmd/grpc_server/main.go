package main

import (
	"log"
	"net"

	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
	pipeline_proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/pipeline"
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
	// Initialize the database
	infrastructure.InitDatabase()

	// Get the database instance
	db := infrastructure.GetDB() // ✅ Fix: Use GetDB() to retrieve *gorm.DB

	// Pass the database instance to the repository adapter
	dbRepo := secondary.NewDatabaseAdapter(db) // ✅ Fix: Passing *gorm.DB correctly

	authService := services.NewAuthService(dbRepo)
	pipelineService := services.NewPipelineService(dbRepo)

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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	grpcServer.GracefulStop()

	log.Println("Servers exited gracefully")
}
