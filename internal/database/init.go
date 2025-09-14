package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// ExecuteSQLFile reads and executes a SQL file
func ExecuteSQLFile(db *sql.DB, filePath string) error {
	log.Printf("📄 Reading SQL file: %s", filePath)
	
	content, err := ioutil.ReadFile(filePath)
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
	
	// First, ensure we're using the correct database
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS alertly DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Printf("⚠️ Could not create database: %v", err)
	}
	
	_, err = db.Exec("USE alertly")
	if err != nil {
		return fmt.Errorf("error selecting database: %v", err)
	}
	
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
		if _, err := ioutil.ReadFile(testSchemaPath); err == nil {
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
func createMinimalSchema(db *sql.DB) error {
	log.Println("🔧 Creating minimal schema...")
	
	queries := []string{
		`CREATE TABLE IF NOT EXISTS incident_categories (
			incca_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(45) NOT NULL,
			icon VARCHAR(45) NULL,
			color VARCHAR(7) NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (incca_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		
		`CREATE TABLE IF NOT EXISTS incident_subcategories (
			incsu_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			incca_id INT UNSIGNED NOT NULL,
			name VARCHAR(100) NOT NULL,
			description TEXT NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (incsu_id),
			KEY fk_subcategory_category_idx (incca_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		
		// Insert basic categories if empty
		`INSERT IGNORE INTO incident_categories (incca_id, name, icon, color) VALUES
			(1, 'Traffic', 'car', '#FF5722'),
			(2, 'Crime', 'shield', '#F44336'),
			(3, 'Fire', 'fire', '#FF9800'),
			(4, 'Medical', 'medical', '#4CAF50'),
			(5, 'Weather', 'cloud', '#2196F3'),
			(6, 'Infrastructure', 'build', '#9C27B0'),
			(7, 'Community', 'people', '#00BCD4'),
			(8, 'Other', 'help', '#607D8B')`,
		
		`INSERT IGNORE INTO incident_subcategories (incca_id, name, description) VALUES
			(1, 'Accident', 'Vehicle collision or accident'),
			(1, 'Traffic Jam', 'Heavy traffic congestion'),
			(2, 'Theft', 'Robbery or theft incident'),
			(3, 'Building Fire', 'Fire in a building'),
			(4, 'Medical Emergency', 'Person needs medical attention'),
			(5, 'Severe Weather', 'Dangerous weather conditions'),
			(6, 'Power Outage', 'Electricity outage'),
			(7, 'Lost Pet', 'Missing pet'),
			(8, 'Other Incident', 'Other type of incident')`,
	}
	
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("⚠️ Error executing query: %v", err)
		}
	}
	
	log.Println("✅ Minimal schema created")
	return nil
}

// CheckAndInitDatabase checks if tables exist and creates them if needed
func CheckAndInitDatabase(db *sql.DB) error {
	// First ensure we're using the right database
	_, err := db.Exec("USE alertly")
	if err != nil {
		log.Printf("⚠️ Database 'alertly' might not exist, creating it...")
		if err := InitializeSchema(db); err != nil {
			return err
		}
		return nil
	}
	
	// Check if incident_categories table exists
	var tableName string
	err = db.QueryRow(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'alertly' 
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