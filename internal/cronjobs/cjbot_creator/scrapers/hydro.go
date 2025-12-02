package scrapers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// Toronto Hydro Outage Map API
	// NOTE: This URL is a PLACEHOLDER - you need to reverse-engineer the actual API endpoint
	// by inspecting network requests on https://www.torontohydro.com/outage-map
	HYDRO_API_URL = "https://api.torontohydro.com/outages/current"
	// Possible alternatives:
	// - ArcGIS REST API endpoint (many utilities use Esri maps)
	// - GeoJSON feed endpoint
)

// HydroScraper handles scraping Toronto Hydro power outages
type HydroScraper struct {
	httpClient *http.Client
}

// NewHydroScraper creates a new Hydro scraper instance
func NewHydroScraper() *HydroScraper {
	return &HydroScraper{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HydroOutageResponse represents the API response structure
// IMPORTANT: This structure is HYPOTHETICAL - update based on actual API response
type HydroOutageResponse struct {
	Outages []HydroOutage `json:"outages"`
}

type HydroOutage struct {
	OutageID          string        `json:"outageId"`
	Status            string        `json:"status"` // "active", "planned", "resolved"
	Cause             string        `json:"cause"`
	CustomersAffected int           `json:"customersAffected"`
	StartTime         string        `json:"startTime"`
	ETR               string        `json:"estimatedRestoration"` // Estimated Time of Restoration
	Location          string        `json:"location"`
	Polygon           [][]float64   `json:"polygon"` // [[lng, lat], [lng, lat], ...]
	Center            *LatLng       `json:"center"`  // May be provided by API
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Scrape fetches active power outages from Toronto Hydro API
func (s *HydroScraper) Scrape() ([]ScrapedIncident, error) {
	log.Printf("⚡ Starting Toronto Hydro scraper...")

	// Fetch JSON from API
	resp, err := s.httpClient.Get(HYDRO_API_URL)
	if err != nil {
		return nil, fmt.Errorf("fetching Hydro API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Hydro API returned status %d", resp.StatusCode)
	}

	// Read and parse JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var apiResponse HydroOutageResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	incidents := []ScrapedIncident{}

	for _, outage := range apiResponse.Outages {
		// Skip resolved outages
		if strings.ToLower(outage.Status) == "resolved" {
			continue
		}

		// Parse timestamps
		startTime, _ := s.parseTimestamp(outage.StartTime)
		etr, _ := s.parseTimestamp(outage.ETR)

		// Convert polygon to our Point format
		polygon := s.convertPolygon(outage.Polygon)

		// Determine coordinates (use API-provided center or calculate centroid)
		var lat, lng *float64
		if outage.Center != nil {
			lat = &outage.Center.Lat
			lng = &outage.Center.Lng
		} else if len(polygon) > 0 {
			// Will be calculated as centroid later by normalizer
			// Leave as nil for now
		}

		incident := ScrapedIncident{
			Source:         "hydro",
			ExternalID:     outage.OutageID,
			RawTitle:       fmt.Sprintf("Power Outage - %d customers affected", outage.CustomersAffected),
			RawDescription: fmt.Sprintf("Cause: %s. Location: %s", outage.Cause, outage.Location),
			RawCategory:    outage.Status, // "active", "planned", etc.
			Address:        outage.Location,
			Latitude:       lat,
			Longitude:      lng,
			Polygon:        polygon,
			Timestamp:      startTime,
			ETR:            &etr,
			Status:         outage.Status,
		}

		incidents = append(incidents, incident)
	}

	log.Printf("✅ Hydro scraper found %d active outages", len(incidents))
	return incidents, nil
}

// convertPolygon converts API polygon format to our Point slice
func (s *HydroScraper) convertPolygon(coords [][]float64) []Point {
	points := []Point{}
	for _, coord := range coords {
		if len(coord) >= 2 {
			points = append(points, Point{
				Lng: coord[0],
				Lat: coord[1],
			})
		}
	}
	return points
}

// parseTimestamp converts Hydro timestamp string to time.Time
func (s *HydroScraper) parseTimestamp(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}

	// Try multiple formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("no matching format for: %s", timeStr)
}

// ScrapeMockData returns mock Hydro data for testing
func (s *HydroScraper) ScrapeMockData() []ScrapedIncident {
	log.Printf("⚠️ Hydro: Using MOCK data (real endpoint not configured)")

	now := time.Now()
	etr1 := now.Add(2 * time.Hour)
	etr2 := now.Add(4 * time.Hour)

	// Mock polygon for downtown Toronto area
	downtownPolygon := []Point{
		{Lat: 43.6500, Lng: -79.3800},
		{Lat: 43.6520, Lng: -79.3800},
		{Lat: 43.6520, Lng: -79.3750},
		{Lat: 43.6500, Lng: -79.3750},
		{Lat: 43.6500, Lng: -79.3800}, // Close polygon
	}

	// Mock polygon for Scarborough area
	scarboroughPolygon := []Point{
		{Lat: 43.7730, Lng: -79.2580},
		{Lat: 43.7750, Lng: -79.2580},
		{Lat: 43.7750, Lng: -79.2530},
		{Lat: 43.7730, Lng: -79.2530},
		{Lat: 43.7730, Lng: -79.2580},
	}

	lat1, lng1 := 43.6510, -79.3775
	lat2, lng2 := 43.7740, -79.2555

	return []ScrapedIncident{
		{
			Source:         "hydro",
			ExternalID:     "hydro_outage_001",
			RawTitle:       "Power Outage - 450 customers affected",
			RawDescription: "Cause: Equipment failure. Location: Downtown Toronto",
			RawCategory:    "unplanned outage",
			Address:        "Downtown Toronto, Bay St & King St area",
			Latitude:       &lat1,
			Longitude:      &lng1,
			Polygon:        downtownPolygon,
			Timestamp:      now.Add(-45 * time.Minute),
			ETR:            &etr1,
			Status:         "active",
		},
		{
			Source:         "hydro",
			ExternalID:     "hydro_outage_002",
			RawTitle:       "Power Outage - 1200 customers affected",
			RawDescription: "Cause: Storm damage. Location: Scarborough",
			RawCategory:    "unplanned outage",
			Address:        "Scarborough, Kennedy Rd & Eglinton Ave area",
			Latitude:       &lat2,
			Longitude:      &lng2,
			Polygon:        scarboroughPolygon,
			Timestamp:      now.Add(-2 * time.Hour),
			ETR:            &etr2,
			Status:         "active",
		},
	}
}
