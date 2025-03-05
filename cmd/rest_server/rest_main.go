package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/streadway/amqp"

	"github.com/sarika-p9/my-pipeline-project/api/http/handlers"
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

	// Initialize REST API handler
	handler := &handlers.PipelineHandler{Service: pipelineService}
	authHandler := &handlers.AuthHandler{Service: authService}
	userHandler := &handlers.UserHandler{Service: authService}

	// Setup Gin router
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "http://localhost:3000" { // ✅ Allow only frontend origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // ✅ Required for credentials
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	r.POST("/register", authHandler.RegisterHandler)
	r.POST("/login", authHandler.LoginHandler)
	r.POST("/logout", authHandler.LogoutHandler)

	// User profile routes
	r.GET("/user/:id", userHandler.GetUserProfile)
	r.PUT("/user/:id", userHandler.UpdateUserProfile)
	r.GET("/pipelines", handler.GetUserPipelines)
	r.GET("/pipelines/:id/stages", handler.GetPipelineStages)
	r.POST("/createpipelines", handler.CreatePipeline)
	r.POST("/pipelines/:id/start", handler.StartPipeline)
	r.GET("/pipelines/:id/status", handler.GetPipelineStatus)
	r.POST("/pipelines/:id/cancel", handler.CancelPipeline)

	// // Create gRPC server
	// grpcServer := grpc.NewServer()
	// authServer := &primary.AuthServer{AuthService: authService}
	// pipelineServer := &primary.PipelineServer{Service: pipelineService}

	// // Register gRPC services
	// proto.RegisterAuthServiceServer(grpcServer, authServer)
	// pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
	// reflection.Register(grpcServer)

	// // Start gRPC server
	// go func() {
	// 	listener, err := net.Listen("tcp", ":50051")
	// 	if err != nil {
	// 		log.Fatalf("Failed to listen on port 50051: %v", err)
	// 	}
	// 	log.Println("Starting gRPC server on port 50051...")
	// 	if err := grpcServer.Serve(listener); err != nil {
	// 		log.Fatalf("Failed to start gRPC server: %v", err)
	// 	}
	// }()

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
	log.Println("Starting API server on port 8080...")
	if err := r.Run(":8080"); err != nil {

		log.Fatalf("Failed to start server: %v", err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	// grpcServer.GracefulStop()
	log.Println("Servers exited gracefully")
}
