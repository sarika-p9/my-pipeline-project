package grpcserver

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sarika-p9/my-pipeline-project/api/grpc/proto"
	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
)

// StartGRPCServer starts the gRPC server
func StartGRPCServer() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pipelineService := services.NewPipelineService(infrastructure.GetDatabaseInstance()) // ✅ Fixed this line
	proto.RegisterPipelineServiceServer(grpcServer, pipelineService)

	// ✅ Enable gRPC reflection
	reflection.Register(grpcServer)

	log.Println("✅ gRPC Server started on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
