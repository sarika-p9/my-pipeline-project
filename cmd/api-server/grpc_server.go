package main

import (
	"log"
	"net"

	pb "github.com/sarikap9/my-pipeline-project/api/grpc/proto"
	"github.com/sarikap9/my-pipeline-project/internal/adapters/primary"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPipelineServiceServer(grpcServer, &primary.PipelineServer{})

	log.Println("gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
