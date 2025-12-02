package scrapers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	// Toronto Police Service - Calls for Service (Public, No Geographic Offense)
	// Updates in real-time, excludes certain sensitive crime types for privacy
	TPS_API_URL = "https://services.arcgis.com/S9th0jAJ7bqgIRjw/arcgis/rest/services/C4S_Public_NoGO/FeatureServer/0/query"
)

// TPSScraper handles scraping Toronto Police Service calls for service
type TPSScraper struct {
	httpClient *http.Client
}

// NewTPSScraper creates a new TPS scraper instance
func NewTPSScraper() *TPSScraper {
	return &TPSScraper{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TPSResponse represents the ArcGIS FeatureServer response
type TPSResponse struct {
	Features []TPSFeature `json:"features"`
}

type TPSFeature struct {
	Attributes TPSAttributes `json:"attributes"`
	Geometry   TPSGeometry   `json:"geometry"`
}

type TPSAttributes struct {
	ObjectID         int64   `json:"OBJECTID"`
	OccurrenceTime   int64   `json:"OCCURRENCE_TIME_AGOL"` // Unix timestamp in milliseconds
	Division         string  `json:"DIVISION"`
	Latitude         float64 `json:"LATITUDE"`
	Longitude        float64 `json:"LONGITUDE"`
	CallTypeCode     string  `json:"CALL_TYPE_CODE"`
	CallType         string  `json:"CALL_TYPE"`
	CrossStreets     string  `json:"CROSS_STREETS"`
}

type TPSGeometry struct {
	X float64 `json:"x"` // Longitude
	Y float64 `json:"y"` // Latitude
}

// Scrape fetches recent police calls for service from TPS API
func (s *TPSScraper) Scrape() ([]ScrapedIncident, error) {
	log.Printf("🚔 Starting TPS scraper...")

	// Build query URL - get last 6 hours, ordered by most recent first
	now := time.Now()
	yesterday := now.Add(-6 * time.Hour)
	// ArcGIS requires date format: "date 'YYYY-MM-DD HH:MM:SS'"
	yesterdayStr := yesterday.Format("2006-01-02 15:04:05")

	// Build WHERE clause
	whereClause := fmt.Sprintf("OCCURRENCE_TIME_AGOL>=date '%s'", yesterdayStr)

	// Build full URL with proper encoding
	baseURL := TPS_API_URL
	params := url.Values{}
	params.Add("f", "json")
	params.Add("resultOffset", "0")
	params.Add("resultRecordCount", "100")
	params.Add("where", whereClause)
	params.Add("orderByFields", "OCCURRENCE_TIME_AGOL DESC")
	params.Add("outFields", "*")
	params.Add("outSR", "4326")

	fullURL := baseURL + "?" + params.Encode()
	log.Printf("🔍 [TPS] Query URL: %s", fullURL)

	// Fetch data
	resp, err := s.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("fetching TPS API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ [TPS] API Error Response: %s", string(body))
		return nil, fmt.Errorf("TPS API returned status %d", resp.StatusCode)
	}

	// Parse JSON
	var response TPSResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	incidents := []ScrapedIncident{}

	// Convert each feature to ScrapedIncident
	for _, feature := range response.Features {
		attr := feature.Attributes

		// Skip if no valid coordinates
		if attr.Latitude == 0 && attr.Longitude == 0 {
			continue
		}

		// Convert timestamp (milliseconds to time.Time)
		timestamp := time.UnixMilli(attr.OccurrenceTime)

		// Build address from cross streets
		address := buildTPSAddress(attr.CrossStreets)

		// Build description with call details
		description := buildTPSDescription(attr.CallType, attr.Division, attr.CrossStreets)

		lat := attr.Latitude
		lng := attr.Longitude

		incident := ScrapedIncident{
			Source:         "tps",
			ExternalID:     fmt.Sprintf("TPS-%d", attr.ObjectID),
			RawTitle:       attr.CallType,
			RawDescription: description,
			RawCategory:    attr.CallTypeCode,
			RawSubcategory: attr.CallType,
			Address:        address,
			Latitude:       &lat,
			Longitude:      &lng,
			Timestamp:      timestamp,
			Status:         "active",
		}

		incidents = append(incidents, incident)
	}

	log.Printf("✅ TPS scraper found %d incidents", len(incidents))
	return incidents, nil
}

// buildTPSAddress creates a formatted address from cross streets
func buildTPSAddress(crossStreets string) string {
	if crossStreets == "" {
		return "Toronto, ON"
	}
	return crossStreets + ", Toronto, ON"
}

// buildTPSDescription creates a detailed description for TPS calls
func buildTPSDescription(callType, division, crossStreets string) string {
	desc := fmt.Sprintf("Call Type: %s", callType)

	if division != "" {
		desc += fmt.Sprintf(" | Division: %s", division)
	}

	if crossStreets != "" {
		desc += fmt.Sprintf(" | Location: %s", crossStreets)
	}

	return desc
}
