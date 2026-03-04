package repository

import (
	"database/sql"

	"web3-tools-backend/internal/model"
)

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(app *model.Application) error {
	return r.db.QueryRow(`
		INSERT INTO applications (job_id, name, email, resume_text, resume_filename, cover_letter)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`,
		app.JobID, app.Name, app.Email, app.ResumeText, app.ResumeFilename, app.CoverLetter,
	).Scan(&app.ID, &app.CreatedAt)
}

func (r *ApplicationRepository) FindByJobID(jobID int) ([]model.Application, error) {
	rows, err := r.db.Query(`
		SELECT id, job_id, name, email, resume_text, resume_filename, cover_letter, status, created_at
		FROM applications WHERE job_id = $1 ORDER BY created_at DESC`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.Application
	for rows.Next() {
		var app model.Application
		err := rows.Scan(
			&app.ID, &app.JobID, &app.Name, &app.Email, &app.ResumeText,
			&app.ResumeFilename, &app.CoverLetter, &app.Status, &app.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *ApplicationRepository) FindByEmail(email string) ([]model.Application, error) {
	rows, err := r.db.Query(`
		SELECT id, job_id, name, email, resume_text, resume_filename, cover_letter, status, created_at
		FROM applications WHERE email = $1 ORDER BY created_at DESC`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.Application
	for rows.Next() {
		var app model.Application
		err := rows.Scan(
			&app.ID, &app.JobID, &app.Name, &app.Email, &app.ResumeText,
			&app.ResumeFilename, &app.CoverLetter, &app.Status, &app.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}
