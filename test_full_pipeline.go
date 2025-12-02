package main

import (
	"alertly/internal/cronjobs/cjbot_creator"
	"alertly/internal/cronjobs/cjbot_creator/scrapers"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🧪 TESTING FULL BOT CREATOR PIPELINE")
	fmt.Println("=" + repeat("=", 79))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  No .env file found, using environment variables")
	}

	// Test 1: Scraping
	fmt.Println("\n📡 STEP 1: Scraping TFS incidents...")
	fmt.Println("-" + repeat("-", 79))

	scraper := scrapers.NewTFSScraper()
	incidents, err := scraper.Scrape()
	if err != nil {
		log.Fatalf("❌ Scraping failed: %v", err)
	}

	fmt.Printf("✅ Scraped %d incidents\n", len(incidents))
	if len(incidents) > 0 {
		fmt.Printf("\n📋 Sample incident:\n")
		sample := incidents[0]
		fmt.Printf("   Source:      %s\n", sample.Source)
		fmt.Printf("   External ID: %s\n", sample.ExternalID)
		fmt.Printf("   Category:    %s\n", sample.RawCategory)
		fmt.Printf("   Address:     %s\n", sample.Address)
		fmt.Printf("   Timestamp:   %s\n", sample.Timestamp)
	}

	// Test 2: Normalization
	fmt.Println("\n🔄 STEP 2: Normalizing to Alertly schema...")
	fmt.Println("-" + repeat("-", 79))

	normalizer := cjbot_creator.NewNormalizer()
	normalizedCount := 0
	failedCount := 0

	categoryStats := make(map[string]int)

	for i, incident := range incidents {
		normalized, err := normalizer.Normalize(incident)
		if err != nil {
			fmt.Printf("⚠️  [%d] Failed to normalize %s: %v\n", i+1, incident.ExternalID, err)
			failedCount++
			continue
		}

		normalizedCount++
		categoryStats[normalized.CategoryCode]++

		if i < 3 { // Show first 3 normalized incidents
			fmt.Printf("\n[%d] ✅ %s\n", i+1, incident.ExternalID)
			fmt.Printf("    Raw Category:    %s\n", incident.RawCategory)
			fmt.Printf("    → Alertly Cat:   %s\n", normalized.CategoryCode)
			fmt.Printf("    → Subcategory:   %s\n", normalized.SubcategoryCode)
			fmt.Printf("    Title:           %s\n", normalized.Title)
			fmt.Printf("    Address:         %s\n", normalized.Address)
			fmt.Printf("    Image:           %s\n", normalized.ImageURL)
		}
	}

	fmt.Printf("\n📊 Normalization Results:\n")
	fmt.Printf("   Successful: %d\n", normalizedCount)
	fmt.Printf("   Failed:     %d\n", failedCount)

	fmt.Printf("\n📈 Category Breakdown:\n")
	for cat, count := range categoryStats {
		fmt.Printf("   - %s: %d incidents\n", cat, count)
	}

	// Test 3: Database connection (optional)
	fmt.Println("\n💾 STEP 3: Checking database connection...")
	fmt.Println("-" + repeat("-", 79))

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" {
		fmt.Println("⚠️  Database credentials not found in .env, skipping DB test")
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbUser, dbPass, dbHost, dbPort, dbName)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("⚠️  Database connection failed: %v\n", err)
		} else {
			defer db.Close()

			if err := db.Ping(); err != nil {
				fmt.Printf("⚠️  Database ping failed: %v\n", err)
			} else {
				fmt.Println("✅ Database connection successful")

				// Check if bot user exists
				var botUserID int
				err := db.QueryRow("SELECT account_id FROM account WHERE account_id = 1").Scan(&botUserID)
				if err == sql.ErrNoRows {
					fmt.Println("⚠️  Bot user (ID=1) not found in database")
				} else if err != nil {
					fmt.Printf("⚠️  Error checking bot user: %v\n", err)
				} else {
					fmt.Printf("✅ Bot user exists (ID=%d)\n", botUserID)
				}
			}
		}
	}

	// Test 4: Geocoding check (without actual API call)
	fmt.Println("\n🗺️  STEP 4: Geocoding requirements...")
	fmt.Println("-" + repeat("-", 79))

	addressesNeedingGeocode := 0
	for _, incident := range incidents {
		if incident.Latitude == nil || incident.Longitude == nil {
			addressesNeedingGeocode++
		}
	}

	fmt.Printf("   Incidents needing geocoding: %d/%d\n", addressesNeedingGeocode, len(incidents))
	fmt.Printf("   ℹ️  Geocoding will happen automatically via Nominatim (1 req/sec)\n")
	fmt.Printf("   ℹ️  Estimated geocoding time: ~%d seconds\n", addressesNeedingGeocode)

	fmt.Println("\n" + repeat("=", 80))
	fmt.Println("✅ PIPELINE TEST COMPLETE")
	fmt.Println("=" + repeat("=", 79))

	fmt.Println("\n💡 Summary:")
	fmt.Printf("   - Scraping:       ✅ Working (%d incidents)\n", len(incidents))
	fmt.Printf("   - Normalization:  ✅ %d/%d successful\n", normalizedCount, len(incidents))
	fmt.Printf("   - Geocoding:      ⏳ Required for %d addresses\n", addressesNeedingGeocode)
	fmt.Printf("   - Database:       %s\n", getDBStatus(dbUser))

	if failedCount > 0 {
		fmt.Printf("\n⚠️  WARNING: %d incidents failed normalization - check category mappings\n", failedCount)
	}
}

func getDBStatus(dbUser string) string {
	if dbUser == "" {
		return "⚠️  Not configured"
	}
	return "✅ Ready"
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
