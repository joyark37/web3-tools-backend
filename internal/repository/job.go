package repository

import (
	"database/sql"

	"web3-tools-backend/internal/model"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) FindAll(category, search string) ([]model.Job, error) {
	query := `
		SELECT id, title, company, email, location, job_type, salary_min, salary_max,
		       category, description, requirements, benefits, tags, status, created_at, updated_at
		FROM jobs WHERE status = 'active'`
	args := []interface{}{}
	argNum := 1

	if category != "" && category != "all" {
		query += " AND category = $1"
		args = append(args, category)
		argNum++
	}

	if search != "" {
		if argNum > 1 {
			query += " AND (title ILIKE $2 OR company ILIKE $2 OR tags ILIKE $2)"
		} else {
			query += " AND (title ILIKE $1 OR company ILIKE $1 OR tags ILIKE $1)"
		}
		args = append(args, "%"+search+"%")
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.Job
	for rows.Next() {
		var job model.Job
		err := rows.Scan(
			&job.ID, &job.Title, &job.Company, &job.Email, &job.Location, &job.Type,
			&job.SalaryMin, &job.SalaryMax, &job.Category, &job.Description,
			&job.Requirements, &job.Benefits, &job.Tags, &job.Status, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *JobRepository) FindByID(id int) (*model.Job, error) {
	var job model.Job
	err := r.db.QueryRow(`
		SELECT id, title, company, email, location, job_type, salary_min, salary_max,
		       category, description, requirements, benefits, tags, status, created_at, updated_at
		FROM jobs WHERE id = $1`, id).Scan(
		&job.ID, &job.Title, &job.Company, &job.Email, &job.Location, &job.Type,
		&job.SalaryMin, &job.SalaryMax, &job.Category, &job.Description,
		&job.Requirements, &job.Benefits, &job.Tags, &job.Status, &job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) Create(job *model.Job) error {
	tags := ""
	if job.Tags != "" {
		tags = job.Tags
	} else {
		tags = job.Title + "," + job.Company
	}

	return r.db.QueryRow(`
		INSERT INTO jobs (title, company, email, location, job_type, salary_min, salary_max, 
		                  category, description, requirements, benefits, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`,
		job.Title, job.Company, job.Email, job.Location, job.Type,
		job.SalaryMin, job.SalaryMax, job.Category, job.Description,
		job.Requirements, job.Benefits, tags,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
}
