//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Open the local HTML file
	file, err := os.Open("internal/cronjobs/cjbot_creator/scrapers/tfs.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🔍 Testing TFS scraper with LOCAL HTML file...")
	fmt.Println("=============================================================================")

	// Find the table
	table := doc.Find("table#chart")
	fmt.Printf("✓ Found table#chart: %v\n", table.Length() > 0)

	tbody := table.Find("tbody")
	fmt.Printf("✓ Found tbody: %v\n", tbody.Length() > 0)

	rows := tbody.Find("tr")
	fmt.Printf("✓ Found %d rows in tbody\n\n", rows.Length())

	// Iterate over rows
	incidentCount := 0
	rows.Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 8 {
			fmt.Printf("⚠️  Row %d: Only %d cells (skipped)\n", i+1, cells.Length())
			return
		}

		primeStreet := strings.TrimSpace(cells.Eq(0).Text())
		crossStreet := strings.TrimSpace(cells.Eq(1).Text())
		timeStr := strings.TrimSpace(cells.Eq(2).Text())
		incidentNumber := strings.TrimSpace(cells.Eq(3).Text())
		incidentType := strings.TrimSpace(cells.Eq(4).Text())
		alarmLevel := strings.TrimSpace(cells.Eq(5).Text())
		station := strings.TrimSpace(cells.Eq(6).Text())

		if incidentNumber == "" {
			fmt.Printf("⚠️  Row %d: No incident number (skipped)\n", i+1)
			return
		}

		incidentCount++
		fmt.Printf("\n[%d] 🔥 Incident: %s\n", incidentCount, incidentNumber)
		fmt.Printf("    Type:         %s\n", incidentType)
		fmt.Printf("    Prime Street: %s\n", primeStreet)
		fmt.Printf("    Cross Street: %s\n", crossStreet)
		fmt.Printf("    Time:         %s\n", timeStr)
		fmt.Printf("    Alarm Level:  %s\n", alarmLevel)
		fmt.Printf("    Station:      %s\n", station)
	})

	fmt.Printf("\n=============================================================================\n")
	fmt.Printf("📊 TOTAL: Found %d valid incidents\n", incidentCount)
}
