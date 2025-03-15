package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	proto "github.com/sarika-p9/my-pipeline-project/api/grpc/proto/pipeline"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Manage pipelines",
}

var createPipelineCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		userID, _ := cmd.Flags().GetString("user")
		pipelineName, _ := cmd.Flags().GetString("pipeline-name")
		stages, _ := cmd.Flags().GetInt("stages")
		isParallel, _ := cmd.Flags().GetBool("parallel")
		stageNames, _ := cmd.Flags().GetString("stage-names")
		if userID == "" || pipelineName == "" || stages <= 0 {
			log.Fatal("❌ User ID, Pipeline Name, and a valid number of stages are required.")
		}

		if _, err := uuid.Parse(userID); err != nil {
			log.Fatal("❌ Invalid user ID format.")
		}
		var stageNamesList []string
		if stageNames != "" {
			stageNamesList = strings.Split(stageNames, ",")
		} else {
			log.Fatal("❌ Stage names are required.")
		}
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := client.CreatePipeline(ctx, &proto.CreatePipelineRequest{
			UserId:       userID,
			PipelineName: pipelineName,
			Stages:       int32(stages),
			IsParallel:   isParallel,
			StageNames:   stageNamesList,
		})
		if err != nil {
			log.Fatalf("❌ Pipeline creation failed: %v", err)
		}

		fmt.Printf("✅ Pipeline created successfully! Pipeline ID: %s\n", resp.PipelineId)
	},
}

var startPipelineCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a pipeline execution",
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID, _ := cmd.Flags().GetString("pipeline-id")
		userID, _ := cmd.Flags().GetString("user-id")
		inputStr, _ := cmd.Flags().GetString("input")
		isParallel, _ := cmd.Flags().GetBool("parallel")
		if pipelineID == "" || userID == "" {
			log.Fatal("❌ Pipeline ID and User ID are required.")
		}
		if _, err := uuid.Parse(userID); err != nil {
			log.Fatal("❌ Invalid user ID format.")
		}
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var inputMap map[string]interface{}
		if inputStr != "" {
			if err := json.Unmarshal([]byte(inputStr), &inputMap); err != nil {
				log.Fatalf("❌ Failed to parse input JSON: %v", err)
			}
		} else {
			inputMap = map[string]interface{}{}
		}
		inputStruct, err := structpb.NewStruct(inputMap)
		if err != nil {
			log.Fatalf("❌ Failed to convert input to protobuf Struct: %v", err)
		}
		inputAny, err := anypb.New(inputStruct)
		if err != nil {
			log.Fatalf("❌ Failed to wrap input in Any: %v", err)
		}
		resp, err := client.StartPipeline(ctx, &proto.StartPipelineRequest{
			PipelineId: pipelineID,
			UserId:     userID,
			Input:      inputAny,
			IsParallel: isParallel,
		})
		if err != nil {
			log.Fatalf("❌ Failed to start pipeline: %v", err)
		}

		fmt.Printf("🚀 Pipeline execution started successfully! Message: %s\n", resp.Message)
	},
}

var cancelPipelineCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a running pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID, _ := cmd.Flags().GetString("pipeline-id")
		userID, _ := cmd.Flags().GetString("user-id")
		isParallel, _ := cmd.Flags().GetBool("parallel")

		if pipelineID == "" || userID == "" {
			log.Fatal("❌ Pipeline ID and User ID are required.")
		}

		if _, err := uuid.Parse(userID); err != nil {
			log.Fatal("❌ Invalid user ID format.")
		}

		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := client.CancelPipeline(ctx, &proto.CancelPipelineRequest{
			PipelineId: pipelineID,
			UserId:     userID,
			IsParallel: isParallel,
		})
		if err != nil {
			log.Fatalf("❌ Failed to cancel pipeline: %v", err)
		}

		fmt.Printf("❌ Pipeline cancelled successfully! Message: %s\n", resp.Message)
	},
}

var getPipelineStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		pipelineID, _ := cmd.Flags().GetString("pipeline-id")
		isParallel, _ := cmd.Flags().GetBool("parallel")
		if pipelineID == "" {
			log.Fatal("❌ Pipeline ID is required.")
		}
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("❌ Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := proto.NewPipelineServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.GetPipelineStatus(ctx, &proto.GetPipelineStatusRequest{
			PipelineId: pipelineID,
			IsParallel: isParallel,
		})
		if err != nil {
			log.Fatalf("❌ Failed to get pipeline status: %v", err)
		}
		fmt.Println("📊 Pipeline Status Report:")
		fmt.Println("------------------------------")
		fmt.Printf("🆔 Pipeline ID: %s\n", pipelineID)
		fmt.Printf("📌 Status: %s\n", resp.Status)
		fmt.Println("------------------------------")
	},
}

func init() {
	pipelineCmd.AddCommand(createPipelineCmd)
	pipelineCmd.AddCommand(startPipelineCmd)
	pipelineCmd.AddCommand(cancelPipelineCmd)
	pipelineCmd.AddCommand(getPipelineStatusCmd)

	createPipelineCmd.Flags().String("user", "", "User ID")
	createPipelineCmd.Flags().String("pipeline-name", "", "Pipeline Name")
	createPipelineCmd.Flags().Int("stages", 0, "Number of stages (required)")
	createPipelineCmd.Flags().Bool("parallel", true, "Parallel execution")
	createPipelineCmd.Flags().String("stage-names", "", "Comma-separated list of stage names")

	createPipelineCmd.MarkFlagRequired("user")
	createPipelineCmd.MarkFlagRequired("pipeline-name")
	createPipelineCmd.MarkFlagRequired("stages")
	createPipelineCmd.MarkFlagRequired("stage-names")

	startPipelineCmd.Flags().String("pipeline-id", "", "Pipeline ID")
	startPipelineCmd.Flags().String("user-id", "", "User ID")
	startPipelineCmd.Flags().String("input", "", "Input for pipeline")
	startPipelineCmd.Flags().Bool("parallel", true, "Run in parallel mode")
	startPipelineCmd.MarkFlagRequired("pipeline-id")
	startPipelineCmd.MarkFlagRequired("user-id")

	getPipelineStatusCmd.Flags().String("pipeline-id", "", "Pipeline ID")
	getPipelineStatusCmd.Flags().Bool("parallel", false, "Check parallel pipeline status")
	getPipelineStatusCmd.MarkFlagRequired("pipeline-id")

}
