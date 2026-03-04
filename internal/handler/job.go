package handler

import (
	"log"
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
		log.Printf("Error listing jobs: %v", err)
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch jobs")
		return
	}

	if jobs == nil {
		jobs = []model.JobResponse{}
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) getJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.service.GetJob(id)
	if err != nil {
		if err.Error() == "job not found" {
			RespondWithError(c, http.StatusNotFound, "Job not found")
			return
		}
		log.Printf("Error getting job %d: %v", id, err)
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch job")
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) createJob(c *gin.Context) {
	var req model.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	// Validate request
	if validationErrs := model.ValidateJobRequest(&req); len(validationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": validationErrs,
		})
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
		log.Printf("Error creating job: %v", err)
		RespondWithError(c, http.StatusInternalServerError, "Failed to create job")
		return
	}

	RespondWithCreated(c, job)
}
