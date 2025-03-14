package primary

import (
	"context"
	"log"

	"github.com/google/uuid"
	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/pipeline"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PipelineServer struct {
	proto.UnimplementedPipelineServiceServer
	Service *services.PipelineService
}

func (s *PipelineServer) CreatePipeline(ctx context.Context, req *proto.CreatePipelineRequest) (*proto.CreatePipelineResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
	}

	pipelineName := req.PipelineName
	if pipelineName == "" {
		pipelineName = "Untitled Pipeline"
	}

	stageNames := make([]string, len(req.StageNames))
	for i, name := range req.StageNames {
		if name == "" {
			stageNames[i] = "Untitled Stage"
		} else {
			stageNames[i] = name
		}
	}

	pipelineID, err := s.Service.CreatePipeline(userID, pipelineName, int(req.Stages), stageNames)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create pipeline: %v", err)
	}

	return &proto.CreatePipelineResponse{PipelineId: pipelineID.String()}, nil
}

func (s *PipelineServer) StartPipeline(ctx context.Context, req *proto.StartPipelineRequest) (*proto.StartPipelineResponse, error) {
	log.Println("[GRPC] Received StartPipeline request...")

	pipelineID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		log.Printf("[ERROR] Invalid pipeline ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid pipeline ID: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Printf("[ERROR] Invalid user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
	}

	if userID == uuid.Nil {
		log.Println("[ERROR] User ID is required but received nil")
		return nil, status.Error(codes.InvalidArgument, "User ID is required")
	}

	var input interface{}
	if req.Input != nil {
		log.Println("[DEBUG] Processing input data...")

		stringWrapper := &wrapperspb.StringValue{}
		if err := req.Input.UnmarshalTo(stringWrapper); err == nil {
			input = stringWrapper.Value
			log.Printf("[DEBUG] Parsed input as string: %v", input)
		} else {
			structValue := &structpb.Struct{}
			if err := req.Input.UnmarshalTo(structValue); err == nil {
				input = structValue.AsMap()
				log.Printf("[DEBUG] Parsed input as JSON object: %v", input)
			} else {
				log.Printf("[ERROR] Failed to unpack input: %v", err)
				return nil, status.Errorf(codes.InvalidArgument, "Invalid input format: %v", err)
			}
		}
	} else {
		log.Println("[DEBUG] No input data provided, using nil")
	}

	go func() {
		log.Printf("[INFO] Starting pipeline execution: %s", pipelineID)
		err := s.Service.StartPipeline(context.Background(), userID, pipelineID, input)
		if err != nil {
			log.Printf("[ERROR] Pipeline execution failed for %s: %v", pipelineID, err)
		} else {
			log.Printf("[INFO] Pipeline execution completed successfully: %s", pipelineID)
		}
	}()

	return &proto.StartPipelineResponse{
		Message: "Pipeline execution started",
	}, nil
}

func (s *PipelineServer) GetPipelineStatus(ctx context.Context, req *proto.GetPipelineStatusRequest) (*proto.GetPipelineStatusResponse, error) {
	pipelineID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid pipeline ID: %v", err)
	}

	stat, err := s.Service.GetPipelineStatus(pipelineID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get pipeline status: %v", err)
	}

	return &proto.GetPipelineStatusResponse{
		PipelineId: pipelineID.String(),
		Status:     stat,
	}, nil
}

func (s *PipelineServer) CancelPipeline(ctx context.Context, req *proto.CancelPipelineRequest) (*proto.CancelPipelineResponse, error) {
	pipelineID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid pipeline ID: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID: %v", err)
	}

	err = s.Service.CancelPipeline(pipelineID, userID)
	if err != nil {
		log.Printf("Error cancelling pipeline %s: %v", pipelineID, err)
		return nil, status.Errorf(codes.Internal, "Failed to cancel pipeline: %v", err)
	}

	return &proto.CancelPipelineResponse{Message: "Pipeline cancelled"}, nil
}
