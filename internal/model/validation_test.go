package model

import (
	"testing"
)

func TestValidateJobRequest_Valid(t *testing.T) {
	minSalary := 100
	maxSalary := 200
	req := &CreateJobRequest{
		Title:        "Senior Engineer",
		Company:      "Tech Corp",
		Email:        "hr@techcorp.com",
		Location:     "Remote",
		Type:         "full-time",
		SalaryMin:    &minSalary,
		SalaryMax:    &maxSalary,
		Category:     "engineering",
		Description:  "Great job opportunity",
		Requirements: "5+ years experience",
		Benefits:     "Competitive salary",
		Tags:         []string{"Go", "React"},
	}

	errs := ValidateJobRequest(req)
	if len(errs) != 0 {
		t.Errorf("Expected no validation errors, got %v", errs)
	}
}

func TestValidateJobRequest_EmptyTitle(t *testing.T) {
	req := &CreateJobRequest{
		Title:   "",
		Company: "Tech Corp",
		Email:   "hr@techcorp.com",
	}

	errs := ValidateJobRequest(req)
	if len(errs) == 0 {
		t.Error("Expected validation error for empty title")
	}
}

func TestValidateJobRequest_InvalidEmail(t *testing.T) {
	req := &CreateJobRequest{
		Title:   "Engineer",
		Company: "Tech Corp",
		Email:   "invalid-email",
	}

	errs := ValidateJobRequest(req)
	if len(errs) == 0 {
		t.Error("Expected validation error for invalid email")
	}
}

func TestValidateJobRequest_InvalidSalaryRange(t *testing.T) {
	minSalary := 200
	maxSalary := 100 // Invalid: min > max
	req := &CreateJobRequest{
		Title:     "Engineer",
		Company:   "Tech Corp",
		Email:     "hr@techcorp.com",
		SalaryMin: &minSalary,
		SalaryMax: &maxSalary,
	}

	errs := ValidateJobRequest(req)
	if len(errs) == 0 {
		t.Error("Expected validation error for invalid salary range")
	}
}

func TestValidateJobRequest_InvalidJobType(t *testing.T) {
	req := &CreateJobRequest{
		Title:   "Engineer",
		Company: "Tech Corp",
		Email:   "hr@techcorp.com",
		Type:    "invalid-type",
	}

	errs := ValidateJobRequest(req)
	if len(errs) == 0 {
		t.Error("Expected validation error for invalid job type")
	}
}

func TestValidateJobRequest_InvalidCategory(t *testing.T) {
	req := &CreateJobRequest{
		Title:    "Engineer",
		Company:  "Tech Corp",
		Email:    "hr@techcorp.com",
		Category: "invalid-category",
	}

	errs := ValidateJobRequest(req)
	if len(errs) == 0 {
		t.Error("Expected validation error for invalid category")
	}
}

func TestValidateApplicationRequest_Valid(t *testing.T) {
	req := &CreateApplicationRequest{
		JobID:         1,
		Name:          "John Doe",
		Email:         "john@example.com",
		ResumeText:    "My resume",
		ResumeFilename: "resume.pdf",
		CoverLetter:   "Cover letter",
	}

	errs := ValidateApplicationRequest(req)
	if len(errs) != 0 {
		t.Errorf("Expected no validation errors, got %v", errs)
	}
}

func TestValidateApplicationRequest_EmptyName(t *testing.T) {
	req := &CreateApplicationRequest{
		JobID:  1,
		Name:   "",
		Email:  "john@example.com",
	}

	errs := ValidateApplicationRequest(req)
	if len(errs) == 0 {
		t.Error("Expected validation error for empty name")
	}
}

func TestValidateApplicationRequest_InvalidJobID(t *testing.T) {
	req := &CreateApplicationRequest{
		JobID: 0, // Invalid
		Name:  "John Doe",
		Email: "john@example.com",
	}

	errs := ValidateApplicationRequest(req)
	if len(errs) == 0 {
		t.Error("Expected validation error for invalid job ID")
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"test.user@example.com", true},
		{"test+tag@example.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tt := range tests {
		result := isValidEmail(tt.email)
		if result != tt.valid {
			t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, result, tt.valid)
		}
	}
}
