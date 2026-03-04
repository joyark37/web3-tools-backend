package model

import (
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var msgs []string
	for _, err := range v {
		msgs = append(msgs, err.Field+": "+err.Message)
	}
	return strings.Join(msgs, ", ")
}

func ValidateJobRequest(req *CreateJobRequest) ValidationErrors {
	var errors ValidationErrors

	// Title validation
	if strings.TrimSpace(req.Title) == "" {
		errors = append(errors, ValidationError{Field: "title", Message: "职位名称不能为空"})
	} else if len(req.Title) > 255 {
		errors = append(errors, ValidationError{Field: "title", Message: "职位名称不能超过255个字符"})
	}

	// Company validation
	if strings.TrimSpace(req.Company) == "" {
		errors = append(errors, ValidationError{Field: "company", Message: "公司名称不能为空"})
	} else if len(req.Company) > 255 {
		errors = append(errors, ValidationError{Field: "company", Message: "公司名称不能超过255个字符"})
	}

	// Email validation
	if strings.TrimSpace(req.Email) == "" {
		errors = append(errors, ValidationError{Field: "email", Message: "邮箱不能为空"})
	} else if !isValidEmail(req.Email) {
		errors = append(errors, ValidationError{Field: "email", Message: "邮箱格式不正确"})
	}

	// Location validation (optional)
	if len(req.Location) > 255 {
		errors = append(errors, ValidationError{Field: "location", Message: "工作地点不能超过255个字符"})
	}

	// Job type validation
	validTypes := map[string]bool{
		"full-time": true,
		"part-time": true,
		"contract":   true,
		"intern":      true,
	}
	if req.Type != "" && !validTypes[req.Type] {
		errors = append(errors, ValidationError{Field: "type", Message: "无效的工作类型"})
	}

	// Category validation
	validCategories := map[string]bool{
		"engineering": true,
		"design":      true,
		"product":     true,
		"marketing":   true,
		"operations": true,
	}
	if req.Category != "" && !validCategories[req.Category] {
		errors = append(errors, ValidationError{Field: "category", Message: "无效的职位类别"})
	}

	// Salary validation
	if req.SalaryMin != nil && req.SalaryMax != nil && *req.SalaryMin > *req.SalaryMax {
		errors = append(errors, ValidationError{Field: "salary", Message: "最低薪资不能大于最高薪资"})
	}

	// Tags validation
	if len(req.Tags) > 20 {
		errors = append(errors, ValidationError{Field: "tags", Message: "标签不能超过20个"})
	}

	return errors
}

func ValidateApplicationRequest(req *CreateApplicationRequest) ValidationErrors {
	var errors ValidationErrors

	// JobID validation
	if req.JobID <= 0 {
		errors = append(errors, ValidationError{Field: "jobId", Message: "无效的职位ID"})
	}

	// Name validation
	if strings.TrimSpace(req.Name) == "" {
		errors = append(errors, ValidationError{Field: "name", Message: "姓名不能为空"})
	} else if len(req.Name) > 255 {
		errors = append(errors, ValidationError{Field: "name", Message: "姓名不能超过255个字符"})
	}

	// Email validation
	if strings.TrimSpace(req.Email) == "" {
		errors = append(errors, ValidationError{Field: "email", Message: "邮箱不能为空"})
	} else if !isValidEmail(req.Email) {
		errors = append(errors, ValidationError{Field: "email", Message: "邮箱格式不正确"})
	}

	// Resume filename validation
	if req.ResumeFilename != "" && len(req.ResumeFilename) > 255 {
		errors = append(errors, ValidationError{Field: "resumeFilename", Message: "简历文件名不能超过255个字符"})
	}

	return errors
}

func isValidEmail(email string) bool {
	// Simple email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
