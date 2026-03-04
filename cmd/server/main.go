package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"web3-tools-backend/internal/handler"
	"web3-tools-backend/internal/repository"
	"web3-tools-backend/internal/service"
)

func main() {
	// Load configuration
	cfg := loadConfig()

	// Connect to database
	dsn := buildDSN(cfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL database:", cfg.DBName)

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize layers
	jobRepo := repository.NewJobRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	
	jobService := service.NewJobService(jobRepo)
	applicationService := service.NewApplicationService(applicationRepo)
	
	jobHandler := handler.NewJobHandler(jobService)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	// Setup Gin router
	router := gin.Default()
	
	// CORS middleware
	router.Use(corsMiddleware())

	// Register routes
	jobHandler.RegisterRoutes(router)
	applicationHandler.RegisterRoutes(router)

	// Start server
	log.Printf("Server starting on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

func loadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "joy"),
		DBPassword: getEnv("DB_PASSWORD", "ServBay.dev"),
		DBName:     getEnv("DB_NAME", "web3_recruitment"),
		Port:       getEnv("PORT", "8080"),
	}
}

func buildDSN(cfg *Config) string {
	return "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS jobs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			company VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			location VARCHAR(255) DEFAULT 'Remote',
			job_type VARCHAR(50) DEFAULT 'full-time',
			salary_min INTEGER,
			salary_max INTEGER,
			category VARCHAR(50) DEFAULT 'engineering',
			description TEXT,
			requirements TEXT,
			benefits TEXT,
			tags TEXT,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS applications (
			id SERIAL PRIMARY KEY,
			job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			resume_text TEXT,
			resume_filename VARCHAR(255),
			cover_letter TEXT,
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_category ON jobs(category)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_applications_job_id ON applications(job_id)`,
		`CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return err
		}
	}

	log.Println("Database migrations completed")
	return nil
}
