package main

import (
	"alertly/internal/cronjob"
	"alertly/internal/database"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// This script can be run as a cronjob to check for expired premium subscriptions
// Usage: go run scripts/premium_expiration_cron.go
// Crontab example: 0 2 * * * /path/to/go run /path/to/premium_expiration_cron.go (daily at 2 AM)

func main() {
	log.Println("🔄 Premium Expiration Cronjob Started")
	startTime := time.Now()

	// Load environment variables
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("Warning: Could not load .env file: %v", err)
		}
	}

	// Initialize database from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	database.InitDB(dsn)
	defer database.DB.Close()

	log.Println("✅ Database connected successfully")

	// Create premium expiration service
	service := cronjob.NewPremiumExpirationService(database.DB)

	// 1. Check and expire premium accounts
	log.Println("🔍 Checking for expired premium accounts...")
	err := service.CheckAndExpirePremiumAccounts()
	if err != nil {
		log.Fatalf("❌ Premium expiration check failed: %v", err)
	}

	// 2. Send expiration warnings (optional - run less frequently)
	log.Println("📧 Checking for expiration warnings...")
	err = service.SendExpirationWarnings()
	if err != nil {
		log.Printf("⚠️ Expiration warnings failed (non-fatal): %v", err)
	}

	// 3. Get and log statistics
	stats, err := service.GetPremiumExpirationStats()
	if err != nil {
		log.Printf("⚠️ Could not get premium stats: %v", err)
	} else {
		log.Printf("📊 Premium Stats: %+v", stats)
	}

	duration := time.Since(startTime)
	log.Printf("✅ Premium Expiration Cronjob Completed in %v", duration)
}
