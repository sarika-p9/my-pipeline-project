package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
	pipeline_proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/pipeline"
	"github.com/sarika-p9/my-pipeline-project/api/http/handlers"
	"github.com/sarika-p9/my-pipeline-project/internal/adapters/primary"
	"github.com/sarika-p9/my-pipeline-project/internal/adapters/secondary"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/messaging"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	rabbitMQConn    *amqp.Connection
	rabbitMQChannel *amqp.Channel
	natsConn        *nats.Conn
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Proceeding with existing environment variables.")
	}
	infrastructure.InitDatabase()
	db := infrastructure.GetDB()
	dbRepo := secondary.NewDatabaseAdapter(db)
	authService := services.NewAuthService(dbRepo)
	pipelineService := services.NewPipelineService(dbRepo)
	handler := &handlers.PipelineHandler{Service: pipelineService}
	authHandler := &handlers.AuthHandler{Service: authService}
	userHandler := &handlers.UserHandler{Service: authService}
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "http://localhost:3000" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
	r.GET("/user/:id", userHandler.GetUserProfile)
	r.PUT("/user/:id", userHandler.UpdateUserProfile)
	r.GET("/pipelines", handler.GetUserPipelines)
	r.GET("/pipelines/:id/stages", handler.GetPipelineStages)
	r.POST("/createpipelines", handler.CreatePipeline)
	r.POST("/pipelines/:id/start", handler.StartPipeline)
	r.GET("/pipelines/:id/status", handler.GetPipelineStatus)
	r.POST("/pipelines/:id/cancel", handler.CancelPipeline)
	grpcServer := grpc.NewServer()
	authServer := &primary.AuthServer{AuthService: authService}
	pipelineServer := &primary.PipelineServer{Service: pipelineService}
	proto.RegisterAuthServiceServer(grpcServer, authServer)
	pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
	reflection.Register(grpcServer)
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
	go func() {
		natsURL := "nats://localhost:4222"
		natsConn = infrastructure.ConnectNATS(natsURL)
		defer natsConn.Close()
	}()
	log.Println("Starting API server on port 8080...")
	if err := r.Run(":8080"); err != nil {

		log.Fatalf("Failed to start server: %v", err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")
	log.Println("Servers exited gracefully")
}
