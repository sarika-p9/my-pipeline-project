package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sarika-p9/my-pipeline-project/internal/services"
)

type PipelineHandler struct {
	Service *services.PipelineService
}

type CreatePipelineRequest struct {
	Name       string   `json:"name"`
	Stages     int      `json:"stages"` // ✅ Change from []models.Stages to int
	IsParallel bool     `json:"is_parallel"`
	UserID     string   `json:"user_id"`
	StageNames []string `json:"stage_names"` // ✅ Ensure this stays as []string
}

// CreatePipeline handles pipeline creation
// CreatePipeline handles pipeline creation
func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	var req CreatePipelineRequest

	// Bind JSON and log the received payload
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ Invalid request payload:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	fmt.Printf("📥 Received CreatePipeline Request: %+v\n", req)

	if req.UserID == "" {
		fmt.Println("❌ Missing User ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// ✅ Validate Stages
	if req.Stages <= 0 {
		fmt.Println("❌ Invalid number of stages:", req.Stages)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of stages! Must be greater than 0."})
		return
	}

	// Convert UserID (string) to uuid.UUID
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		fmt.Println("❌ Invalid User ID format:", req.UserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// ✅ Log details before creating pipeline
	fmt.Printf("🛠️ Creating Pipeline: Name=%s, Stages=%d, UserID=%s, StageNames=%v\n", req.Name, req.Stages, userUUID, req.StageNames)

	// ✅ Call service to create pipeline (stages inserted inside service)
	pipelineID, err := h.Service.CreatePipeline(userUUID, req.Name, req.Stages, req.StageNames)
	if err != nil {
		fmt.Println("❌ Failed to create pipeline:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pipeline"})
		return
	}

	// ✅ Log success
	fmt.Println("✅ Pipeline and Stages Created Successfully, ID:", pipelineID)

	// Return response
	c.JSON(http.StatusAccepted, gin.H{
		"message":     "Pipeline and stages created",
		"pipeline_id": pipelineID.String(),
	})
}

type StartPipelineRequest struct {
	Input      interface{} `json:"input"`
	IsParallel bool        `json:"is_parallel"`
	UserID     string      `json:"user_id"` // Change from uuid.UUID to string
}

func (h *PipelineHandler) StartPipeline(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id")) // Ensure pipelineID is converted to uuid.UUID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var req StartPipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Convert UserID from string to uuid.UUID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// ✅ Only start the pipeline now, do NOT insert stages
	go func() {
		h.Service.StartPipeline(context.Background(), userID, pipelineID, req.Input)
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Pipeline execution started", "pipeline_id": pipelineID})
}

type GetPipelineStatusRequest struct {
	IsParallel bool `json:"is_parallel"`
}

// GetPipelineStatus retrieves the current status of a pipeline
func (h *PipelineHandler) GetPipelineStatus(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	// var req GetPipelineStatusRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	// 	return
	// }

	// Get the "is_parallel" query parameter (defaults to false if not provided)

	// status, err := h.Service.GetPipelineStatus(pipelineID, req.IsParallel)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
	// 	return
	// }
	status, err := h.Service.GetPipelineStatus(pipelineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pipeline_id": pipelineID, "status": status})
}

type CancelPipelineRequest struct {
	IsParallel bool   `json:"is_parallel"`
	UserID     string `json:"user_id"`
}

// CancelPipeline cancels an ongoing pipeline execution
func (h *PipelineHandler) CancelPipeline(c *gin.Context) {
	// Convert pipelineID from string to uuid.UUID
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var req CancelPipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Convert req.UserID from string to uuid.UUID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Pass converted uuid.UUID values to CancelPipeline service
	err = h.Service.CancelPipeline(pipelineID, userID)
	if err != nil {
		log.Printf("Error cancelling pipeline: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel pipeline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline cancelled", "pipeline_id": pipelineID})
}

func (h *PipelineHandler) GetUserPipelines(c *gin.Context) {
	userID := c.Query("user_id") // Fetch user_id from query parameters
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	pipelines, err := h.Service.GetPipelinesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipelines"})
		return
	}

	c.JSON(http.StatusOK, pipelines)
}

// GetPipelineStages fetches the stages of a pipeline
func (h *PipelineHandler) GetPipelineStages(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	stages, err := h.Service.GetPipelineStages(pipelineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipeline stages"})
		return
	}

	c.JSON(http.StatusOK, stages)
}
