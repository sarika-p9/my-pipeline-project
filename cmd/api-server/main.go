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
	"github.com/sarikap9/my-pipeline-project/api/grpc/proto"
	"github.com/sarikap9/my-pipeline-project/api/http/routes"
	"github.com/sarikap9/my-pipeline-project/internal/grpcserver"
	"github.com/sarikap9/my-pipeline-project/internal/infrastructure"
	"github.com/sarikap9/my-pipeline-project/internal/messaging"
	"github.com/sarikap9/my-pipeline-project/internal/middleware"
	"github.com/sarikap9/my-pipeline-project/internal/services"
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

// initRabbitMQ initializes the RabbitMQ connection.
func initRabbitMQ() {
	var err error
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	rabbitMQConn, rabbitMQChannel, err = messaging.ConnectRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("✅ Connected to RabbitMQ")
}

// closeRabbitMQ closes the RabbitMQ connection.
func closeRabbitMQ() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
	if rabbitMQConn != nil {
		rabbitMQConn.Close()
	}
}

// initNATS initializes the NATS connection.
func initNATS() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222" // Default for local
	}
	natsConn = infrastructure.ConnectNATS(natsURL)
}

// closeNATS closes the NATS connection.
func closeNATS() {
	if natsConn != nil {
		natsConn.Close()
	}
}

// ConsumePipelineEvents listens to the "pipelines" queue and processes messages.
func ConsumePipelineEvents(handler func(string) error) error {
	queue, err := rabbitMQChannel.QueueDeclare(
		"pipelines",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := rabbitMQChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("📥 Received a pipeline event: %s", d.Body)
			err := handler(string(d.Body))
			if err != nil {
				log.Printf("❌ Error processing message: %v", err)
				d.Nack(false, true)
				continue
			}
			d.Ack(false)
		}
	}()

	return nil
}

func main() {
	// Load environment variables.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Supabase.
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set in the .env file")
	}
	db, err := infrastructure.InitSupabaseWithGORM(supabaseURL, supabaseKey)
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	log.Println("DB connection and migration successful:", db)

	// Initialize gRPC, RabbitMQ, and NATS.
	initGRPC()
	defer closeGRPC()

	initRabbitMQ()
	defer closeRabbitMQ()

	initNATS()
	defer closeNATS()

	// Subscribe to pipeline updates.
	services.SubscribeToEvents(natsConn, "pipeline.updates")

	// Start the gRPC server.
	go grpcserver.StartGRPCServer()

	// Start consuming pipeline events.
	if err := ConsumePipelineEvents(func(msg string) error {
		log.Printf("Processing pipeline event: %s", msg)
		services.PublishRealTimeEvent(natsConn, "pipeline.updates", "Processed: "+msg)
		return nil
	}); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Initialize messaging producer.
	producer, err := messaging.NewProducer(rabbitMQChannel, "pipeline_tasks")
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	// Publish a test message.
	if err := producer.Publish("Initial test message"); err != nil {
		log.Printf("Error publishing initial message: %v", err)
	}

	// Initialize Gin router.
	r := gin.Default()

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

	// API route: Create new pipeline.
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

		services.PublishRealTimeEvent(natsConn, "pipeline.updates", "Pipeline Created: "+newPipeline.Name)
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

	infrastructure.InitDB()

	// Setup router
	r = routes.SetupRouter()

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	// Apply AuthMiddleware to secure routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	// Protected Route Example
	authorized.POST("/pipelines", func(c *gin.Context) {
		var newPipeline struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&newPipeline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Get userID from context (set by AuthMiddleware)
		userID, _ := c.Get("userID")

		// Create pipeline logic here using userID
		c.JSON(http.StatusCreated, gin.H{"message": "Pipeline created", "user_id": userID, "pipeline_name": newPipeline.Name})
	})

	// Run API server.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not run API server: %v", err)
		}
	}()

	// Graceful shutdown.
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
