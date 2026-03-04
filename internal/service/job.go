package service

import (
	"strings"

	"web3-tools-backend/internal/model"
	"web3-tools-backend/internal/repository"
)

type JobService struct {
	repo *repository.JobRepository
}

func NewJobService(repo *repository.JobRepository) *JobService {
	return &JobService{repo: repo}
}

func (s *JobService) ListJobs(category, search string) ([]model.JobResponse, error) {
	jobs, err := s.repo.FindAll(category, search)
	if err != nil {
		return nil, err
	}

	responses := make([]model.JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = s.toResponse(&job)
	}

	return responses, nil
}

func (s *JobService) GetJob(id int) (*model.JobResponse, error) {
	job, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(job)
	return &resp, nil
}

func (s *JobService) CreateJob(req *model.CreateJobRequest) (*model.JobResponse, error) {
	tags := ""
	if len(req.Tags) > 0 {
		tags = strings.Join(req.Tags, ",")
	}

	job := &model.Job{
		Title:        req.Title,
		Company:      req.Company,
		Email:        req.Email,
		Location:     req.Location,
		Type:         req.Type,
		SalaryMin:    req.SalaryMin,
		SalaryMax:    req.SalaryMax,
		Category:     req.Category,
		Description:  req.Description,
		Requirements: req.Requirements,
		Benefits:     req.Benefits,
		Tags:         tags,
		Status:       "active",
	}

	if err := s.repo.Create(job); err != nil {
		return nil, err
	}

	resp := s.toResponse(job)
	return &resp, nil
}

func (s *JobService) toResponse(job *model.Job) model.JobResponse {
	var tags []string
	if job.Tags != "" {
		tags = strings.Split(job.Tags, ",")
	}

	return model.JobResponse{
		ID:           job.ID,
		Title:        job.Title,
		Company:      job.Company,
		Email:        job.Email,
		Location:     job.Location,
		Type:         job.Type,
		SalaryMin:    job.SalaryMin,
		SalaryMax:    job.SalaryMax,
		Category:     job.Category,
		Description:  job.Description,
		Requirements: job.Requirements,
		Benefits:     job.Benefits,
		Tags:         tags,
		Status:       job.Status,
		PostedAt:     job.CreatedAt.Format("2006-01-02"),
	}
}

// Simple itoa implementation
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}
