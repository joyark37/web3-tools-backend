package service

import (
	"testing"

	"web3-tools-backend/internal/model"
)

// MockJobRepository is a mock implementation for testing
type MockJobRepository struct {
	jobs   map[int]*model.Job
	nextID int
}

func NewMockJobRepository() *MockJobRepository {
	return &MockJobRepository{
		jobs:   make(map[int]*model.Job),
		nextID: 1,
	}
}

func (m *MockJobRepository) FindAll(category, search string) ([]model.Job, error) {
	var result []model.Job
	for _, job := range m.jobs {
		if category != "" && category != "all" && job.Category != category {
			continue
		}
		if search != "" {
			// Simple search check
			if !contains(job.Title, search) && !contains(job.Company, search) {
				continue
			}
		}
		result = append(result, *job)
	}
	return result, nil
}

func (m *MockJobRepository) FindByID(id int) (*model.Job, error) {
	if job, ok := m.jobs[id]; ok {
		return job, nil
	}
	return nil, nil
}

func (m *MockJobRepository) Create(job *model.Job) error {
	job.ID = m.nextID
	m.jobs[m.nextID] = job
	m.nextID++
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// JobServiceWithMock creates a job service with mock repository for testing
type JobServiceWithMock struct {
	repo *MockJobRepository
}

func NewJobServiceWithMock() *JobServiceWithMock {
	return &JobServiceWithMock{
		repo: NewMockJobRepository(),
	}
}

func (s *JobServiceWithMock) ListJobs(category, search string) ([]model.JobResponse, error) {
	jobs, err := s.repo.FindAll(category, search)
	if err != nil {
		return nil, err
	}

	responses := make([]model.JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = model.JobResponse{
			ID:        job.ID,
			Title:     job.Title,
			Company:   job.Company,
			Email:     job.Email,
			Location:  job.Location,
			Type:      job.Type,
			Category:  job.Category,
			Status:    job.Status,
			PostedAt:  job.CreatedAt.Format("2006-01-02"),
		}
	}

	return responses, nil
}

func (s *JobServiceWithMock) CreateJob(req *model.CreateJobRequest) (*model.JobResponse, error) {
	job := &model.Job{
		Title:     req.Title,
		Company:   req.Company,
		Email:     req.Email,
		Location:  req.Location,
		Type:      req.Type,
		Category:  req.Category,
		Status:    "active",
	}

	if err := s.repo.Create(job); err != nil {
		return nil, err
	}

	return &model.JobResponse{
		ID:        job.ID,
		Title:     job.Title,
		Company:   job.Company,
		Email:     job.Email,
		Location:  job.Location,
		Type:      job.Type,
		Category:  job.Category,
		Status:    job.Status,
		PostedAt:  job.CreatedAt.Format("2006-01-02"),
	}, nil
}

func TestJobService_CreateJob(t *testing.T) {
	svc := NewJobServiceWithMock()

	req := &model.CreateJobRequest{
		Title:    "Senior Engineer",
		Company:  "Tech Corp",
		Email:    "hr@techcorp.com",
		Location: "Remote",
		Type:     "full-time",
		Category: "engineering",
	}

	job, err := svc.CreateJob(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if job.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, job.Title)
	}

	if job.ID == 0 {
		t.Error("Expected job ID to be set")
	}
}

func TestJobService_ListJobs(t *testing.T) {
	svc := NewJobServiceWithMock()

	// Create some jobs
	jobs := []*model.CreateJobRequest{
		{Title: "Engineer", Company: "Corp A", Category: "engineering"},
		{Title: "Designer", Company: "Corp B", Category: "design"},
		{Title: "Manager", Company: "Corp C", Category: "engineering"},
	}

	for _, j := range jobs {
		svc.CreateJob(j)
	}

	// Test listing all
	allJobs, err := svc.ListJobs("", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(allJobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(allJobs))
	}

	// Test filtering by category
	engineeringJobs, err := svc.ListJobs("engineering", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(engineeringJobs) != 2 {
		t.Errorf("Expected 2 engineering jobs, got %d", len(engineeringJobs))
	}

	// Test searching
	searchResults, err := svc.ListJobs("", "Engineer")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(searchResults) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(searchResults))
	}
}
