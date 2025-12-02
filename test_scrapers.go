package main

import (
	"alertly/internal/cronjobs/cjbot_creator/scrapers"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_scrapers.go [tfs|hydro|all]")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run test_scrapers.go tfs     # Test Toronto Fire Services scraper")
		fmt.Println("  go run test_scrapers.go hydro   # Test Toronto Hydro scraper")
		fmt.Println("  go run test_scrapers.go all     # Test all scrapers")
		os.Exit(1)
	}

	scraperType := os.Args[1]

	switch scraperType {
	case "tfs":
		testTFS()
	case "hydro":
		testHydro()
	case "all":
		testTFS()
		fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
		testHydro()
	default:
		log.Printf("❌ Unknown scraper type: %s", scraperType)
		os.Exit(1)
	}
}

// ============================================
// TFS SCRAPER TEST
// ============================================

func testTFS() {
	fmt.Println("🔥 TESTING TFS SCRAPER (Toronto Fire Services)")
	fmt.Println(strings.Repeat("=", 80))

	scraper := scrapers.NewTFSScraper()

	// Try real scraping
	fmt.Println("\n📡 Attempting real scraping from TFS website...")
	incidents, err := scraper.Scrape()

	if err != nil {
		fmt.Printf("⚠️  Real scraping failed: %v\n", err)
		fmt.Println("\n📦 Falling back to MOCK data...")
		incidents = scraper.ScrapeMockData()
	} else {
		fmt.Printf("✅ Real scraping successful!\n")
	}

	// Display results
	fmt.Printf("\n📊 RESULTS: Found %d incidents\n", len(incidents))
	fmt.Println(strings.Repeat("-", 80))

	if len(incidents) == 0 {
		fmt.Println("❌ No incidents found")
		return
	}

	for i, incident := range incidents {
		fmt.Printf("\n[%d] 🔥 TFS Incident\n", i+1)
		fmt.Printf("    External ID:  %s\n", incident.ExternalID)
		fmt.Printf("    Title:        %s\n", incident.RawTitle)
		fmt.Printf("    Category:     %s\n", incident.RawCategory)
		fmt.Printf("    Address:      %s\n", incident.Address)
		fmt.Printf("    Timestamp:    %s\n", incident.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Status:       %s\n", incident.Status)

		if incident.Latitude != nil && incident.Longitude != nil {
			fmt.Printf("    Coordinates:  (%.6f, %.6f)\n", *incident.Latitude, *incident.Longitude)
		} else {
			fmt.Printf("    Coordinates:  ⚠️  Not available (needs geocoding)\n")
		}

		if incident.RawDescription != "" {
			fmt.Printf("    Description:  %s\n", incident.RawDescription)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}

// ============================================
// HYDRO SCRAPER TEST
// ============================================

func testHydro() {
	fmt.Println("⚡ TESTING HYDRO SCRAPER (Toronto Hydro Power Outages)")
	fmt.Println(strings.Repeat("=", 80))

	scraper := scrapers.NewHydroScraper()

	// Try real scraping
	fmt.Println("\n📡 Attempting real scraping from Hydro API...")
	incidents, err := scraper.Scrape()

	if err != nil {
		fmt.Printf("⚠️  Real scraping failed: %v\n", err)
		fmt.Println("\n📦 Falling back to MOCK data...")
		incidents = scraper.ScrapeMockData()
	} else {
		fmt.Printf("✅ Real scraping successful!\n")
	}

	// Display results
	fmt.Printf("\n📊 RESULTS: Found %d outages\n", len(incidents))
	fmt.Println(strings.Repeat("-", 80))

	if len(incidents) == 0 {
		fmt.Println("❌ No outages found")
		return
	}

	for i, incident := range incidents {
		fmt.Printf("\n[%d] ⚡ Power Outage\n", i+1)
		fmt.Printf("    External ID:  %s\n", incident.ExternalID)
		fmt.Printf("    Title:        %s\n", incident.RawTitle)
		fmt.Printf("    Category:     %s\n", incident.RawCategory)
		fmt.Printf("    Address:      %s\n", incident.Address)
		fmt.Printf("    Timestamp:    %s\n", incident.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Status:       %s\n", incident.Status)

		if incident.Latitude != nil && incident.Longitude != nil {
			fmt.Printf("    Coordinates:  (%.6f, %.6f)\n", *incident.Latitude, *incident.Longitude)
		} else {
			fmt.Printf("    Coordinates:  ⚠️  Not available\n")
		}

		if incident.ETR != nil {
			fmt.Printf("    ETR:          %s\n", incident.ETR.Format("2006-01-02 15:04:05"))
		}

		if len(incident.Polygon) > 0 {
			fmt.Printf("    Polygon:      %d points\n", len(incident.Polygon))
			fmt.Printf("                  First point: (%.6f, %.6f)\n",
				incident.Polygon[0].Lat, incident.Polygon[0].Lng)
		}

		if incident.RawDescription != "" {
			fmt.Printf("    Description:  %s\n", incident.RawDescription)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}
