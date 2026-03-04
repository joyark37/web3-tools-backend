package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"web3-tools-backend/internal/model"
	"web3-tools-backend/internal/service"
)

type ApplicationHandler struct {
	service *service.ApplicationService
}

func NewApplicationHandler(s *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: s}
}

func (h *ApplicationHandler) RegisterRoutes(r *gin.Engine) {
	apps := r.Group("/api/applications")
	{
		apps.POST("", h.createApplication)
		apps.GET("/job/:jobId", h.getApplicationsByJob)
	}
}

func (h *ApplicationHandler) createApplication(c *gin.Context) {
	var req model.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	// Validate request
	if validationErrs := model.ValidateApplicationRequest(&req); len(validationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": validationErrs,
		})
		return
	}

	app, err := h.service.CreateApplication(&req)
	if err != nil {
		log.Printf("Error creating application: %v", err)
		RespondWithError(c, http.StatusInternalServerError, "Failed to submit application")
		return
	}

	RespondWithCreated(c, app)
}

func (h *ApplicationHandler) getApplicationsByJob(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("jobId"))
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid job ID")
		return
	}

	apps, err := h.service.GetApplicationsByJob(jobID)
	if err != nil {
		log.Printf("Error getting applications for job %d: %v", jobID, err)
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch applications")
		return
	}

	if apps == nil {
		apps = []model.ApplicationResponse{}
	}

	c.JSON(http.StatusOK, apps)
}
