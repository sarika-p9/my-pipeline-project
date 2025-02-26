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
	"github.com/sarika-p9/my-pipeline-project/internal/controllers"
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
// initRabbitMQ initializes the RabbitMQ connection.
func initRabbitMQ() {
	var err error
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	rabbitMQConn, rabbitMQChannel, err = messaging.ConnectRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	// ❌ Remove this redundant log
	// log.Println("✅ Connected to RabbitMQ")
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
		natsURL = "nats://localhost:4222"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set in the .env file")
	}

	// ✅ Initialize DB and set it globally
	db, err := infrastructure.InitSupabaseWithGORM(supabaseURL, supabaseKey)
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	infrastructure.DB = db // ✅ Ensure global DB is set
	log.Println("✅ DB connection and migration successful!")

	initGRPC()
	defer closeGRPC()

	initRabbitMQ()
	defer closeRabbitMQ()

	initNATS()
	defer closeNATS()

	services.SubscribeToEvents(natsConn, "pipeline.updates")

	go grpcserver.StartGRPCServer()

	if err := ConsumePipelineEvents(func(msg string) error {
		log.Printf("Processing pipeline event: %s", msg)
		services.PublishRealTimeEvent(natsConn, "pipeline.updates", "Processed: "+msg)
		return nil
	}); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	producer, err := messaging.NewProducer(rabbitMQChannel, "pipeline_tasks")
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	if err := producer.Publish("Initial test message"); err != nil {
		log.Printf("Error publishing initial message: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	authorized.POST("/pipelines", func(c *gin.Context) {
		var newPipeline struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&newPipeline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		userID, _ := c.Get("userID")
		c.JSON(http.StatusCreated, gin.H{"message": "Pipeline created", "user_id": userID, "pipeline_name": newPipeline.Name})
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API Server is running!"})
	})

	r.GET("/pipelines", func(c *gin.Context) {
		resp, err := grpcClient.ListPipelines(context.Background(), &proto.ListPipelinesRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pipelines": resp.GetPipelines()})
	})

	r.GET("/workers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"workers": []string{"Worker1", "Worker2"}})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not run API server: %v", err)
		}
	}()

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
