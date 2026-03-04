package handler

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.service.CreateApplication(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, app)
}

func (h *ApplicationHandler) getApplicationsByJob(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("jobId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	apps, err := h.service.GetApplicationsByJob(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apps)
}
