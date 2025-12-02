package scrapers

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// TFS uses a live XML feed that's loaded via JavaScript
	TFS_XML_URL = "https://www.toronto.ca/data/fire/livecad.xml"
	// The HTML page loads data from this XML via AJAX
)

// TFSActiveIncidents represents the XML root element
type TFSActiveIncidents struct {
	XMLName    xml.Name   `xml:"tfs_active_incidents"`
	UpdateTime string     `xml:"update_from_db_time"`
	Events     []TFSEvent `xml:"event"`
}

// TFSEvent represents a single incident event
type TFSEvent struct {
	PrimeStreet   string `xml:"prime_street"`
	CrossStreets  string `xml:"cross_streets"`
	DispatchTime  string `xml:"dispatch_time"`
	EventNum      string `xml:"event_num"`
	EventType     string `xml:"event_type"`
	AlarmLevel    string `xml:"alarm_lev"`
	Beat          string `xml:"beat"`
	UnitsDispatched string `xml:"units_disp"`
}

// TFSScraper handles scraping Toronto Fire Services active incidents
type TFSScraper struct {
	httpClient *http.Client
}

// NewTFSScraper creates a new TFS scraper instance
func NewTFSScraper() *TFSScraper {
	return &TFSScraper{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Scrape fetches active fire incidents from TFS XML feed
func (s *TFSScraper) Scrape() ([]ScrapedIncident, error) {
	log.Printf("🔥 Starting TFS scraper...")

	// Fetch XML feed
	resp, err := s.httpClient.Get(TFS_XML_URL)
	if err != nil {
		return nil, fmt.Errorf("fetching TFS XML: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TFS returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Parse XML
	var activeIncidents TFSActiveIncidents
	if err := xml.Unmarshal(body, &activeIncidents); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	incidents := []ScrapedIncident{}

	// Convert each XML event to ScrapedIncident
	for _, event := range activeIncidents.Events {
		// Skip if no incident number (invalid event)
		if event.EventNum == "" {
			continue
		}

		// Parse timestamp (format: 2025-11-25T15:23:44)
		timestamp, err := s.parseTimestamp(event.DispatchTime)
		if err != nil {
			log.Printf("⚠️ TFS: Failed to parse timestamp '%s': %v", event.DispatchTime, err)
			timestamp = time.Now() // Fallback to current time
		}

		// Build address string
		address := buildAddress(event.PrimeStreet, event.CrossStreets)

		// Build description with additional info
		description := fmt.Sprintf("Incident: %s", event.EventType)
		if event.AlarmLevel != "0" {
			description += fmt.Sprintf(" | Alarm Level: %s", event.AlarmLevel)
		}
		if event.Beat != "" {
			description += fmt.Sprintf(" | Station: %s", event.Beat)
		}
		if event.UnitsDispatched != "" {
			description += fmt.Sprintf(" | Units: %s", event.UnitsDispatched)
		}

		incident := ScrapedIncident{
			Source:         "tfs",
			ExternalID:     event.EventNum, // Use official incident number
			RawTitle:       event.EventType,
			RawDescription: description,
			RawCategory:    event.EventType,
			RawSubcategory: "",
			Address:        address,
			Latitude:       nil, // Will be geocoded later
			Longitude:      nil,
			Timestamp:      timestamp,
			Status:         "active",
		}

		incidents = append(incidents, incident)
	}

	log.Printf("✅ TFS scraper found %d incidents", len(incidents))
	return incidents, nil
}

// parseTimestamp converts TFS timestamp string to time.Time
func (s *TFSScraper) parseTimestamp(timeStr string) (time.Time, error) {
	// TFS uses ISO 8601 format: 2025-11-25T15:23:44
	formats := []string{
		time.RFC3339,                // 2025-11-25T15:23:44Z
		"2006-01-02T15:04:05",       // 2025-11-25T15:23:44 (TFS actual format)
		"2006-01-02T15:04:05-07:00", // With timezone
		"2006-01-02 15:04:05",       // Fallback
		"2006-01-02 15:04",          // Fallback
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("no matching format for: %s", timeStr)
}

// buildAddress constructs a complete address from prime street and cross street
func buildAddress(primeStreet, crossStreet string) string {
	// Clean up empty or placeholder values
	primeStreet = strings.TrimSpace(primeStreet)
	crossStreet = strings.TrimSpace(crossStreet)

	// Handle postal code only cases (e.g., "M6S", "M3A")
	if len(primeStreet) == 3 && primeStreet[0] == 'M' {
		return primeStreet + ", Toronto, ON"
	}

	// Build full address
	if primeStreet == "" || primeStreet == " " {
		if crossStreet != "" && crossStreet != " " && crossStreet != "/" {
			return crossStreet + ", Toronto, ON"
		}
		return "Toronto, ON" // Fallback
	}

	// If we have both streets
	if crossStreet != "" && crossStreet != " " && crossStreet != "/" {
		return primeStreet + " near " + crossStreet + ", Toronto, ON"
	}

	// Just prime street
	return primeStreet + ", Toronto, ON"
}

// ScrapeMockData returns mock TFS data for testing when real API is unavailable
func (s *TFSScraper) ScrapeMockData() []ScrapedIncident {
	log.Printf("⚠️ TFS: Using MOCK data (real endpoint not configured)")

	return []ScrapedIncident{
		{
			Source:         "tfs",
			ExternalID:     "tfs_mock_001",
			RawTitle:       "Structure Fire",
			RawDescription: "Fire crews responding to structure fire",
			RawCategory:    "STRUCTURE FIRE",
			Address:        "100 Queen St W, Toronto",
			Timestamp:      time.Now().Add(-10 * time.Minute),
			Status:         "active",
		},
		{
			Source:         "tfs",
			ExternalID:     "tfs_mock_002",
			RawTitle:       "Medical Emergency",
			RawDescription: "Fire paramedics dispatched",
			RawCategory:    "MEDICAL CALL",
			Address:        "200 Yonge St, Toronto",
			Timestamp:      time.Now().Add(-25 * time.Minute),
			Status:         "active",
		},
		{
			Source:         "tfs",
			ExternalID:     "tfs_mock_003",
			RawTitle:       "Vehicle Fire",
			RawDescription: "Car fire on highway",
			RawCategory:    "VEHICLE FIRE",
			Address:        "Gardiner Expressway & Spadina Ave, Toronto",
			Timestamp:      time.Now().Add(-5 * time.Minute),
			Status:         "active",
		},
	}
}

// sanitizeForID removes special characters from string for use in IDs
func sanitizeForID(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	return s
}
