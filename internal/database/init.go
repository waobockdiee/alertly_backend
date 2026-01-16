package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ExecuteSQLFile reads and executes a SQL file
func ExecuteSQLFile(db *sql.DB, filePath string) error {
	log.Printf("📄 Reading SQL file: %s", filePath)
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	
	// Split by semicolon but handle multi-line statements
	sqlContent := string(content)
	
	// Remove MySQL specific comments and settings that might cause issues
	lines := strings.Split(sqlContent, "\n")
	var cleanedLines []string
	for _, line := range lines {
		// Skip MySQL Workbench specific comments and SET statements
		if strings.HasPrefix(strings.TrimSpace(line), "--") ||
			strings.HasPrefix(strings.TrimSpace(line), "SET @OLD_") ||
			strings.HasPrefix(strings.TrimSpace(line), "SET SQL_MODE") ||
			strings.HasPrefix(strings.TrimSpace(line), "SET @OLD_SQL_MODE") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}
	
	cleanedSQL := strings.Join(cleanedLines, "\n")
	
	// Split into individual statements
	statements := strings.Split(cleanedSQL, ";")
	
	successCount := 0
	errorCount := 0
	
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		// Add semicolon back
		stmt = stmt + ";"
		
		// Execute the statement
		_, err := db.Exec(stmt)
		if err != nil {
			// Log error but continue with other statements
			if !strings.Contains(err.Error(), "already exists") &&
				!strings.Contains(err.Error(), "Duplicate entry") {
				log.Printf("⚠️ Statement %d failed: %v", i+1, err)
				errorCount++
			}
		} else {
			successCount++
		}
	}
	
	log.Printf("✅ Executed %d statements successfully, %d errors", successCount, errorCount)
	return nil
}

// InitializeSchema creates the necessary tables and initial data from SQL files
func InitializeSchema(db *sql.DB) error {
	log.Println("🔧 Initializing database schema from SQL files...")

	// PostgreSQL: Database is already selected via connection string
	// No need for CREATE DATABASE or USE statements
	// The connection is already established to the correct database

	// Verify connection is working
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("error verifying database connection: %v", err)
	}
	log.Println("✅ Database connection verified")
	
	// Define the paths to SQL files
	// These paths are relative to where the binary runs
	basePath := "assets/db"
	
	// Try different possible locations for the SQL files
	possiblePaths := []string{
		basePath,
		filepath.Join("..", "..", basePath),
		filepath.Join("/var/task", basePath),
		filepath.Join("/app", basePath),
	}
	
	var schemaPath, insertsPath string
	for _, path := range possiblePaths {
		testSchemaPath := filepath.Join(path, "db.sql")
		testInsertsPath := filepath.Join(path, "inserts.sql")
		
		// Check if files exist at this path
		if _, err := os.ReadFile(testSchemaPath); err == nil {
			schemaPath = testSchemaPath
			insertsPath = testInsertsPath
			break
		}
	}
	
	if schemaPath == "" {
		// If files not found, create minimal tables to get started
		log.Println("⚠️ SQL files not found, creating minimal schema...")
		return createMinimalSchema(db)
	}
	
	// Execute schema file
	log.Println("📦 Executing schema file...")
	if err := ExecuteSQLFile(db, schemaPath); err != nil {
		log.Printf("⚠️ Error executing schema: %v", err)
		// Try minimal schema as fallback
		return createMinimalSchema(db)
	}
	
	// Check if we need to insert initial data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM incident_categories").Scan(&count)
	if err != nil {
		log.Printf("⚠️ Error checking categories: %v", err)
	}
	
	if count == 0 {
		// Execute inserts file
		log.Println("📝 Executing inserts file...")
		if err := ExecuteSQLFile(db, insertsPath); err != nil {
			log.Printf("⚠️ Error executing inserts: %v", err)
		}
	} else {
		log.Printf("ℹ️ Database already has %d categories, skipping inserts", count)
	}
	
	return nil
}

// createMinimalSchema creates just the essential tables to get the app running
// PostgreSQL version
func createMinimalSchema(db *sql.DB) error {
	log.Println("🔧 Creating minimal schema (PostgreSQL)...")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS incident_categories (
			incca_id SERIAL PRIMARY KEY,
			name VARCHAR(45) NOT NULL,
			icon VARCHAR(45) NULL,
			color VARCHAR(7) NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS incident_subcategories (
			incsu_id SERIAL PRIMARY KEY,
			incca_id INTEGER NOT NULL,
			name VARCHAR(100) NOT NULL,
			description TEXT NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		// Create index for subcategories
		`CREATE INDEX IF NOT EXISTS fk_subcategory_category_idx ON incident_subcategories (incca_id)`,

		// Insert basic categories if empty (ON CONFLICT DO NOTHING = PostgreSQL's INSERT IGNORE)
		`INSERT INTO incident_categories (incca_id, name, icon, color) VALUES
			(1, 'Traffic', 'car', '#FF5722'),
			(2, 'Crime', 'shield', '#F44336'),
			(3, 'Fire', 'fire', '#FF9800'),
			(4, 'Medical', 'medical', '#4CAF50'),
			(5, 'Weather', 'cloud', '#2196F3'),
			(6, 'Infrastructure', 'build', '#9C27B0'),
			(7, 'Community', 'people', '#00BCD4'),
			(8, 'Other', 'help', '#607D8B')
		ON CONFLICT (incca_id) DO NOTHING`,

		`INSERT INTO incident_subcategories (incca_id, name, description) VALUES
			(1, 'Accident', 'Vehicle collision or accident'),
			(1, 'Traffic Jam', 'Heavy traffic congestion'),
			(2, 'Theft', 'Robbery or theft incident'),
			(3, 'Building Fire', 'Fire in a building'),
			(4, 'Medical Emergency', 'Person needs medical attention'),
			(5, 'Severe Weather', 'Dangerous weather conditions'),
			(6, 'Power Outage', 'Electricity outage'),
			(7, 'Lost Pet', 'Missing pet'),
			(8, 'Other Incident', 'Other type of incident')
		ON CONFLICT DO NOTHING`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			// Log but continue - tables may already exist with data
			log.Printf("⚠️ Error executing query (may be normal): %v", err)
		}
	}

	log.Println("✅ Minimal schema created")
	return nil
}

// CheckAndInitDatabase checks if tables exist and creates them if needed
// PostgreSQL version
func CheckAndInitDatabase(db *sql.DB) error {
	// PostgreSQL: Database is already selected via connection string
	// No need for USE statement

	// Check if incident_categories table exists in the public schema
	var tableName string
	err := db.QueryRow(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'incident_categories'
		LIMIT 1
	`).Scan(&tableName)

	if err == sql.ErrNoRows {
		log.Println("📦 Tables not found, initializing database...")
		return InitializeSchema(db)
	} else if err != nil {
		return fmt.Errorf("error checking tables: %v", err)
	}

	// Check if tables have data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM incident_categories").Scan(&count)
	if err != nil || count == 0 {
		log.Println("📦 Tables exist but are empty, initializing data...")
		return InitializeSchema(db)
	}

	log.Printf("✅ Database ready with %d categories", count)
	return nil
}