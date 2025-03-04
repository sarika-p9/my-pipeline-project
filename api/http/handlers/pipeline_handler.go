package handlers

import (
	"context"
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
	Stages int       `json:"stages"`
	UserID uuid.UUID `json:"user_id"` // Extracted from the request
}

// CreatePipeline handles pipeline creation
func (h *PipelineHandler) CreatePipeline(c *gin.Context) {
	var req CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	pipelineID, err := h.Service.CreatePipeline(req.UserID, req.Stages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pipeline"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Pipeline created", "pipeline_id": pipelineID})
}

type StartPipelineRequest struct {
	Input  interface{} `json:"input"`
	UserID uuid.UUID   `json:"user_id"`
}

// StartPipeline handles pipeline execution
func (h *PipelineHandler) StartPipeline(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	var req StartPipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	go func() {
		h.Service.StartPipeline(context.Background(), req.UserID, pipelineID, req.Input)
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Pipeline execution started", "pipeline_id": pipelineID})
}

// GetPipelineStatus retrieves the current status of a pipeline
func (h *PipelineHandler) GetPipelineStatus(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	status, err := h.Service.GetPipelineStatus(pipelineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pipeline_id": pipelineID, "status": status})
}

type CancelPipelineRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

// CancelPipeline cancels an ongoing pipeline execution
func (h *PipelineHandler) CancelPipeline(c *gin.Context) {
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

	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	err = h.Service.CancelPipeline(pipelineID, req.UserID)
	if err != nil {
		log.Printf("Error cancelling pipeline: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel pipeline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline cancelled", "pipeline_id": pipelineID})
}

// GetPipelineStatusHandler handles the request to fetch pipeline status
func (h *PipelineHandler) GetPipelineStatusHandler(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	status, err := h.Service.GetPipelineStatus(pipelineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pipeline_id": pipelineID.String(),
		"status":      status,
	})
}
