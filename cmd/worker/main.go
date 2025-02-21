package main

import (
	"context"
	"log"
	"net"

	"github.com/sarikap9/my-pipeline-project/api/grpc/proto" // Adjust path if needed
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedPipelineServiceServer
}

// Implement the CreatePipeline function
func (s *server) CreatePipeline(ctx context.Context, req *proto.CreatePipelineRequest) (*proto.CreatePipelineResponse, error) {
	log.Println("Received request to create pipeline:", req.Name)
	return &proto.CreatePipelineResponse{Success: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterPipelineServiceServer(grpcServer, &server{})

	log.Println("gRPC server started on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
