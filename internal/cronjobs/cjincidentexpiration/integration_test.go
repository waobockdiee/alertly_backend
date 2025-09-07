package cjincidentexpiration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory: %v", err)
	}
	log.Printf("Test running from working directory: %s", wd)

	envPath := filepath.Join(".env")
	err = godotenv.Load(envPath)
	if err != nil {
		envPath = filepath.Join("../../..", ".env")
		err = godotenv.Load(envPath)
	}

	if err != nil {
		log.Fatalf("FATAL: Could not load .env file. Error: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func TestRepository_Integration(t *testing.T) {
	exec := func(query string, args ...interface{}) sql.Result {
		result, err := db.Exec(query, args...)
		if err != nil {
			t.Fatalf("Failed to execute query: %s, error: %v", query, err)
		}
		return result
	}

	// 1. Setup
	winnerUser := exec("INSERT INTO account (email, password, nickname, score, credibility) VALUES (?, ?, ?, ?, ?)", "winner-repo@test.com", "pass", "winner-repo", 100, 7.0)
	winnerID, _ := winnerUser.LastInsertId()
	defer exec("DELETE FROM account WHERE account_id = ?", winnerID)

	category := exec("INSERT INTO incident_categories (name, code) VALUES (?, ?)", "Test Cat Repo", "test_cat_repo")
	catID, _ := category.LastInsertId()
	defer exec("DELETE FROM incident_categories WHERE inca_id = ?", catID)

	subcategory := exec("INSERT INTO incident_subcategories (inca_id, name, code, default_duration_hours) VALUES (?, ?, ?, ?)", catID, "Test Sub Repo", "test_sub_repo", 1)
	subcatID, _ := subcategory.LastInsertId()
	defer exec("DELETE FROM incident_subcategories WHERE insu_id = ?", subcatID)

	expiredTime := time.Now().Add(-48 * time.Hour)
	cluster := exec("INSERT INTO incident_clusters (insu_id, credibility, is_active, created_at, account_id) VALUES (?, ?, ?, ?, ?)", subcatID, 8.0, "1", expiredTime, winnerID)
	clusterID, _ := cluster.LastInsertId()
	defer exec("DELETE FROM incident_clusters WHERE incl_id = ?", clusterID)

	// Sanity check: ensure the cluster was inserted before proceeding.
	var checkID int64
	err := db.QueryRow("SELECT incl_id FROM incident_clusters WHERE incl_id = ?", clusterID).Scan(&checkID)
	if err != nil {
		t.Fatalf("SANITY CHECK FAILED: Could not select test cluster right after inserting it. DB Error: %v", err)
	}
	if checkID != clusterID {
		t.Fatalf("SANITY CHECK FAILED: Inserted cluster ID %d but selected back %d", clusterID, checkID)
	}

	// Deep diagnostic: select the raw data back and print it.
	var ( 
		read_is_active string
		read_created_at time.Time
		time_comparison_result bool
	)
	diagQuery := "SELECT is_active, created_at, (NOW() >= TIMESTAMPADD(HOUR, 1, created_at)) AS time_check FROM incident_clusters WHERE incl_id = ?"
	err = db.QueryRow(diagQuery, clusterID).Scan(&read_is_active, &read_created_at, &time_comparison_result)
	if err != nil {
		t.Fatalf("DEEP DIAGNOSTIC FAILED: %v", err)
	}
	t.Logf("DEEP DIAGNOSTIC FOR CLUSTER %d: is_active='%s', created_at='%s', time_comparison_evaluates_to='%t'", 
		clusterID, read_is_active, read_created_at, time_comparison_result)

	repo := NewRepository(db)

	// 2. Test GetExpiredClusters
	expiredClusters, err := repo.GetExpiredClusters()
	if err != nil {
		t.Fatalf("GetExpiredClusters failed: %v", err)
	}

	foundTestCluster := false
	for _, c := range expiredClusters {
		if c.ID == clusterID {
			foundTestCluster = true
			if !c.Credibility.Valid || c.Credibility.Float64 != 8.0 {
				t.Errorf("Expected credibility for test cluster to be 8.0, but got %v", c.Credibility)
			}
			break
		}
	}
	if !foundTestCluster {
		t.Fatalf("Test cluster with ID %d was not found in expired clusters list", clusterID)
	}

	// 3. Test UpdateUserStats
	err = repo.UpdateUserStats(winnerID, 15.0, 1.5)
	if err != nil {
		t.Fatalf("UpdateUserStats failed: %v", err)
	}
	var finalWinnerScore, finalWinnerCred float64
	err = db.QueryRow("SELECT score, credibility FROM account WHERE account_id = ?", winnerID).Scan(&finalWinnerScore, &finalWinnerCred)
	if err != nil {
		t.Fatalf("Failed to query winner's final stats: %v", err)
	}
	if finalWinnerScore != 115.0 || finalWinnerCred != 8.5 {
		t.Errorf("Expected winner stats to be (115.0, 8.5), but got (%.1f, %.1f)", finalWinnerScore, finalWinnerCred)
	}

	// 4. Test MarkClusterProcessed
	err = repo.MarkClusterProcessed(clusterID)
	if err != nil {
		t.Fatalf("MarkClusterProcessed failed: %v", err)
	}
	var isActive string
	err = db.QueryRow("SELECT is_active FROM incident_clusters WHERE incl_id = ?", clusterID).Scan(&isActive)
	if err != nil {
		t.Fatalf("Failed to query cluster final status: %v", err)
	}
	if isActive != "0" {
		t.Errorf("Expected cluster to be inactive ('0'), but got '%s'", isActive)
	}
}
