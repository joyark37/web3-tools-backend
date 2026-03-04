package model

import "time"

type Application struct {
	ID            int        `json:"id"`
	JobID         int        `json:"jobId"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	ResumeText    string     `json:"resumeText,omitempty"`
	ResumeFilename string    `json:"resumeFilename,omitempty"`
	CoverLetter   string     `json:"coverLetter,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type ApplicationResponse struct {
	ID              int    `json:"id"`
	JobID           int    `json:"jobId"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ResumeFilename  string `json:"resumeFilename,omitempty"`
	CoverLetter    string `json:"coverLetter,omitempty"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"`
}

type CreateApplicationRequest struct {
	JobID         int    `json:"jobId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	ResumeText    string `json:"resumeText"`
	ResumeFilename string `json:"resumeFilename"`
	CoverLetter   string `json:"coverLetter"`
}
