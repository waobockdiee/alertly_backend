//go:build ignore

package main

import (
	"alertly/internal/cronjobs/cjbot_creator"
	"alertly/internal/database"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🧪 TESTING TPS SCRAPER")
	fmt.Println(strings.Repeat("=", 80))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  No .env file found, using environment variables")
	}

	// Connect to database
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbPort == "" {
		dbPort = "3306"
	}

	fmt.Printf("Connecting to database: %s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)

	// Initialize database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbUser, dbPass, dbHost, dbPort, dbName)
	database.InitDB(dsn)

	fmt.Println("\n📊 BEFORE EXECUTION:")
	fmt.Println(strings.Repeat("-", 80))

	// Run TPS cronjob
	fmt.Println("\n🔥 Running TPS cronjob manually...")
	repo := cjbot_creator.NewRepository(database.DB)
	svc := cjbot_creator.NewService(repo)
	svc.RunTPS()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ TPS CRONJOB TEST COMPLETED")
	fmt.Println(strings.Repeat("=", 80))
}
