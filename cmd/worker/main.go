package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sarika-p9/my-pipeline-project/api/grpc/proto"
)

// PipelineServiceServer implements the gRPC service
type PipelineServiceServer struct {
	proto.UnimplementedPipelineServiceServer
}

// CreatePipeline creates a new pipeline with a unique ID
func (s *PipelineServiceServer) CreatePipeline(ctx context.Context, req *proto.CreatePipelineRequest) (*proto.CreatePipelineResponse, error) {
	pipelineID := uuid.New().String()

	pipeline := models.Pipeline{
		ID:        pipelineID,
		Name:      req.GetName(),
		UserID:    req.GetUserId(),
		Status:    "processing",
		CreatedAt: time.Now(),
	}

	// ✅ Store pipeline in Supabase
	err := infrastructure.InsertPipeline(&pipeline)
	if err != nil {
		log.Printf("Failed to insert pipeline: %v", err)
		return nil, err
	}

	log.Printf("Created Pipeline: ID=%s, Name=%s, UserID=%s", pipelineID, req.GetName(), req.GetUserId())

	// Simulate async processing
	go func(pipelineID string) {
		time.Sleep(5 * time.Second) // Simulate work
		infrastructure.UpdatePipelineStatus(pipelineID, "completed")
		log.Printf("Pipeline processed: ID=%s", pipelineID)
	}(pipelineID)

	return &proto.CreatePipelineResponse{
		Id:      pipelineID,
		Message: "Pipeline created successfully and is being processed",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pipelineService := &PipelineServiceServer{}

	// Register the service
	proto.RegisterPipelineServiceServer(grpcServer, pipelineService)

	// Enable reflection
	reflection.Register(grpcServer)

	log.Println("gRPC server started on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
