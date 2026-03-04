package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"web3-tools-backend/internal/model"
	"web3-tools-backend/internal/service"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(s *service.JobService) *JobHandler {
	return &JobHandler{service: s}
}

func (h *JobHandler) RegisterRoutes(r *gin.Engine) {
	jobs := r.Group("/api/jobs")
	{
		jobs.GET("", h.listJobs)
		jobs.POST("", h.createJob)
		jobs.GET("/:id", h.getJob)
	}
}

func (h *JobHandler) listJobs(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")

	jobs, err := h.service.ListJobs(category, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) getJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	job, err := h.service.GetJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) createJob(c *gin.Context) {
	var req model.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.Location == "" {
		req.Location = "Remote"
	}
	if req.Type == "" {
		req.Type = "full-time"
	}
	if req.Category == "" {
		req.Category = "engineering"
	}

	job, err := h.service.CreateJob(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, job)
}
