package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	cfg := loadConfig()
	dsn := buildDSN(cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database, running migrations...")

	// Get migrations directory
	migrationsDir := getEnv("MIGRATIONS_DIR", "./migrations")
	
	// Run migrations
	if err := runMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully!")
}

func runMigrations(db *sql.DB, dir string) error {
	// Read migration files
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file, err)
		}

		// Split by semicolon and execute each statement
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}

			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute migration from %s: %v\nStatement: %s", file, err, stmt)
			}
		}

		log.Printf("Applied migration: %s", filepath.Base(file))
	}

	return nil
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func loadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "joy"),
		DBPassword: getEnv("DB_PASSWORD", "ServBay.dev"),
		DBName:     getEnv("DB_NAME", "web3_recruitment"),
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
