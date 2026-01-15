//go:build ignore

package main

import (
	"alertly/internal/cronjobs/cjbot_creator"
	"alertly/internal/database"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🧪 TESTING TFS DEDUPLICATION")
	fmt.Println(repeat("=", 80))

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
	database.Init(dbUser, dbPass, dbHost, dbPort, dbName)

	fmt.Println("\n📊 BEFORE EXECUTION:")
	fmt.Println(repeat("-", 80))

	// Run TFS cronjob
	fmt.Println("\n🔥 Running TFS cronjob manually...")
	repo := cjbot_creator.NewRepository(database.DB)
	svc := cjbot_creator.NewService(repo)
	processedCount := svc.RunTFS()

	fmt.Println("\n" + repeat("=", 80))
	fmt.Printf("✅ TFS CRONJOB COMPLETED - Processed: %d incidents\n", processedCount)
	fmt.Println(repeat("=", 80))
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
