package service

import (
	"web3-tools-backend/internal/model"
	"web3-tools-backend/internal/repository"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
}

func NewApplicationService(repo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (s *ApplicationService) CreateApplication(req *model.CreateApplicationRequest) (*model.ApplicationResponse, error) {
	app := &model.Application{
		JobID:         req.JobID,
		Name:          req.Name,
		Email:         req.Email,
		ResumeText:    req.ResumeText,
		ResumeFilename: req.ResumeFilename,
		CoverLetter:   req.CoverLetter,
		Status:        "pending",
	}

	if err := s.repo.Create(app); err != nil {
		return nil, err
	}

	resp := &model.ApplicationResponse{
		ID:             app.ID,
		JobID:          app.JobID,
		Name:           app.Name,
		Email:          app.Email,
		ResumeFilename: app.ResumeFilename,
		CoverLetter:   app.CoverLetter,
		Status:         app.Status,
		CreatedAt:      app.CreatedAt.Format("2006-01-02"),
	}

	return resp, nil
}

func (s *ApplicationService) GetApplicationsByJob(jobID int) ([]model.ApplicationResponse, error) {
	apps, err := s.repo.FindByJobID(jobID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ApplicationResponse, len(apps))
	for i, app := range apps {
		responses[i] = model.ApplicationResponse{
			ID:             app.ID,
			JobID:          app.JobID,
			Name:           app.Name,
			Email:          app.Email,
			ResumeFilename: app.ResumeFilename,
			CoverLetter:   app.CoverLetter,
			Status:         app.Status,
			CreatedAt:      app.CreatedAt.Format("2006-01-02"),
		}
	}

	return responses, nil
}
