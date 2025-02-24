package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sarikap9/my-pipeline-project/api/grpc/proto"
)

// PipelineServiceServer implements the gRPC service
type PipelineServiceServer struct {
	proto.UnimplementedPipelineServiceServer
	pipelines map[string]string // Stores pipeline ID and name
}

// CreatePipeline creates a new pipeline with a unique ID
func (s *PipelineServiceServer) CreatePipeline(ctx context.Context, req *proto.CreatePipelineRequest) (*proto.CreatePipelineResponse, error) {
	id := uuid.New().String()
	s.pipelines[id] = req.GetName()

	log.Printf("Created Pipeline: ID=%s, Name=%s", id, req.GetName())

	// Simulate async processing
	go func(pipelineID, pipelineName string) {
		log.Printf("Processing pipeline: ID=%s, Name=%s", pipelineID, pipelineName)
		time.Sleep(5 * time.Second) // Simulate work
		log.Printf("Pipeline processed: ID=%s, Name=%s", pipelineID, pipelineName)
	}(id, req.GetName())

	return &proto.CreatePipelineResponse{
		Id:      id,
		Message: "Pipeline created successfully and is being processed",
	}, nil
}

// ListPipelines returns all created pipelines
func (s *PipelineServiceServer) ListPipelines(ctx context.Context, req *proto.ListPipelinesRequest) (*proto.ListPipelinesResponse, error) {
	var pipelineNames []string
	for _, name := range s.pipelines {
		pipelineNames = append(pipelineNames, name)
	}

	return &proto.ListPipelinesResponse{
		Pipelines: pipelineNames,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pipelineService := &PipelineServiceServer{pipelines: make(map[string]string)}

	// Register the service
	proto.RegisterPipelineServiceServer(grpcServer, pipelineService)

	// Enable reflection
	reflection.Register(grpcServer)

	log.Println("gRPC server started on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
