package model

import "time"

type Job struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Company      string     `json:"company"`
	Email        string     `json:"email"`
	Location     string     `json:"location"`
	Type         string     `json:"type"`
	SalaryMin    *int       `json:"salaryMin,omitempty"`
	SalaryMax    *int       `json:"salaryMax,omitempty"`
	Category     string     `json:"category"`
	Description  string     `json:"description"`
	Requirements string     `json:"requirements"`
	Benefits     string     `json:"benefits"`
	Tags         string     `json:"tags"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type JobResponse struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Company      string   `json:"company"`
	Email        string   `json:"email"`
	Location     string   `json:"location"`
	Type         string   `json:"type"`
	SalaryMin    *int     `json:"salaryMin,omitempty"`
	SalaryMax    *int     `json:"salaryMax,omitempty"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Requirements string   `json:"requirements"`
	Benefits     string   `json:"benefits"`
	Tags         []string `json:"tags"`
	Status       string   `json:"status"`
	PostedAt     string   `json:"postedAt"`
}

type CreateJobRequest struct {
	Title        string   `json:"title" binding:"required"`
	Company      string   `json:"company" binding:"required"`
	Email        string   `json:"email" binding:"required,email"`
	Location     string   `json:"location"`
	Type         string   `json:"type"`
	SalaryMin    *int     `json:"salaryMin"`
	SalaryMax    *int     `json:"salaryMax"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Requirements string   `json:"requirements"`
	Benefits     string   `json:"benefits"`
	Tags         []string `json:"tags"`
}
