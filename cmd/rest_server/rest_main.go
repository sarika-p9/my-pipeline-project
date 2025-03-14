// package main

// import (
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// 	"github.com/nats-io/nats.go"
// 	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
// 	pipeline_proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/pipeline"
// 	"github.com/sarika-p9/my-pipeline-project/api/http/handlers"
// 	"github.com/sarika-p9/my-pipeline-project/internal/adapters/primary"
// 	"github.com/sarika-p9/my-pipeline-project/internal/adapters/secondary"
// 	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
// 	"github.com/sarika-p9/my-pipeline-project/internal/messaging"
// 	"github.com/sarika-p9/my-pipeline-project/internal/services"
// 	"github.com/streadway/amqp"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// var (
// 	rabbitMQConn    *amqp.Connection
// 	rabbitMQChannel *amqp.Channel
// 	natsConn        *nats.Conn
// )

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("Warning: No .env file found. Proceeding with existing environment variables.")
// 	}
// 	infrastructure.InitDatabase()
// 	db := infrastructure.GetDB()
// 	dbRepo := secondary.NewDatabaseAdapter(db)
// 	authService := services.NewAuthService(dbRepo)
// 	pipelineService := services.NewPipelineService(dbRepo)
// 	handler := &handlers.PipelineHandler{Service: pipelineService}
// 	authHandler := &handlers.AuthHandler{Service: authService}
// 	userHandler := &handlers.UserHandler{Service: authService}
// 	r := gin.Default()
// 	r.Use(func(c *gin.Context) {
// 		origin := c.Request.Header.Get("Origin")
// 		if origin == "http://localhost:3000" {
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		}
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusOK)
// 			return
// 		}
// 		c.Next()
// 	})
// 	r.POST("/register", authHandler.RegisterHandler)
// 	r.POST("/login", authHandler.LoginHandler)
// 	r.POST("/logout", authHandler.LogoutHandler)
// 	r.GET("/user/:id", userHandler.GetUserProfile)
// 	r.PUT("/user/:id", userHandler.UpdateUserProfile)
// 	r.GET("/pipelines", handler.GetUserPipelines)
// 	r.GET("/pipelines/:id/stages", handler.GetPipelineStages)
// 	r.POST("/createpipelines", handler.CreatePipeline)
// 	r.POST("/pipelines/:id/start", handler.StartPipeline)
// 	r.GET("/pipelines/:id/status", handler.GetPipelineStatus)
// 	r.POST("/pipelines/:id/cancel", handler.CancelPipeline)
// 	grpcServer := grpc.NewServer()
// 	authServer := &primary.AuthServer{AuthService: authService}
// 	pipelineServer := &primary.PipelineServer{Service: pipelineService}
// 	proto.RegisterAuthServiceServer(grpcServer, authServer)
// 	pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
// 	reflection.Register(grpcServer)
// 	go func() {
// 		listener, err := net.Listen("tcp", ":50051")
// 		if err != nil {
// 			log.Fatalf("Failed to listen on port 50051: %v", err)
// 		}
// 		log.Println("Starting gRPC server on port 50051...")
// 		if err := grpcServer.Serve(listener); err != nil {
// 			log.Fatalf("Failed to start gRPC server: %v", err)
// 		}
// 	}()
// 	go func() {
// 		rabbitURL := "amqp://guest:guest@localhost:5672/"
// 		var err error
// 		rabbitMQConn, rabbitMQChannel, err = messaging.ConnectRabbitMQ(rabbitURL)
// 		if err != nil {
// 			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 		}
// 		defer rabbitMQChannel.Close()
// 		defer rabbitMQConn.Close()
// 	}()
// 	go func() {
// 		natsURL := "nats://localhost:4222"
// 		natsConn = infrastructure.ConnectNATS(natsURL)
// 		defer natsConn.Close()
// 	}()
// 	log.Println("Starting API server on port 8080...")
// 	if err := r.Run(":8080"); err != nil {

// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
// 	<-quit
// 	log.Println("Shutting down servers...")
// 	log.Println("Servers exited gracefully")
// }

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/authentication"
	pipeline_proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/pipeline"
	"github.com/sarika-p9/my-pipeline-project/api/http/handlers"
	"github.com/sarika-p9/my-pipeline-project/internal/adapters/primary"
	"github.com/sarika-p9/my-pipeline-project/internal/adapters/secondary"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/messaging"
	"github.com/sarika-p9/my-pipeline-project/internal/middleware"
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
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func RESTServer(authService *services.AuthService, pipelineService *services.PipelineService, wg *sync.WaitGroup) {
	defer wg.Done()
	authMiddleware := middleware.AuthMiddleware()
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
	r.POST("/register", gin.WrapF(authHandler.RegisterHandler))
	r.POST("/login", gin.WrapF(authHandler.LoginHandler))
	r.POST("/logout", authHandler.LogoutHandler)
	r.GET("/user/:id", authMiddleware, userHandler.GetUserProfile)
	r.PUT("/user/:id", authMiddleware, userHandler.UpdateUserProfile)
	r.GET("/pipelines", authMiddleware, handler.GetUserPipelines)
	r.GET("/pipelines/:id/stages", authMiddleware, handler.GetPipelineStages)
	r.POST("/createpipelines", authMiddleware, handler.CreatePipeline)
	r.POST("/pipelines/:id/start", authMiddleware, handler.StartPipeline)
	r.GET("/pipelines/:id/status", authMiddleware, handler.GetPipelineStatus)
	r.POST("/pipelines/:id/cancel", authMiddleware, handler.CancelPipeline)
	r.DELETE("/api/pipelines/:pipelineID", authHandler.DeletePipelineHandler)
	r.GET("/ws", func(c *gin.Context) {
		infrastructure.WebSocket.HandleConnections(c)
	})
	go infrastructure.WebSocket.StartBroadcaster()

	log.Println("Starting REST API on port 8080...")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func GRPCServer(authService *services.AuthService, pipelineService *services.PipelineService, wg *sync.WaitGroup) {
	defer wg.Done()
	grpcServer := grpc.NewServer()
	authServer := &primary.AuthServer{AuthService: authService}
	pipelineServer := &primary.PipelineServer{Service: pipelineService}
	proto.RegisterAuthServiceServer(grpcServer, authServer)
	pipeline_proto.RegisterPipelineServiceServer(grpcServer, pipelineServer)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}
	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

// func startFrontendServer(wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	fs := http.FileServer(http.Dir("cmd/"))
// 	http.Handle("/", fs)

// 	log.Println("Starting Frontend server on port 3000...")
// 	if err := http.ListenAndServe(":3000", nil); err != nil {
// 		log.Fatalf("Failed to start Frontend server: %v", err)
// 	}
// }

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Proceeding with existing environment variables.")
	}
	infrastructure.InitDatabase()
	db := infrastructure.GetDB()
	dbRepo := secondary.NewDatabaseAdapter(db)
	authService := services.NewAuthService(dbRepo)
	pipelineService := services.NewPipelineService(dbRepo)

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

	var wg sync.WaitGroup
	wg.Add(3)
	go RESTServer(authService, pipelineService, &wg)
	go GRPCServer(authService, pipelineService, &wg)
	//go startFrontendServer(&wg)
	wg.Wait()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")
	log.Println("Servers exited gracefully")
}
