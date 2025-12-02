package cjbot_creator

import "time"

// ScrapedIncident and Point types are defined in scrapers/types.go
// Import them from there to avoid circular dependencies

// NormalizedIncident represents an incident mapped to Alertly's schema
type NormalizedIncident struct {
	Title          string
	Description    string
	CategoryCode   string  // Alertly category (e.g., "crime", "fire_incident")
	SubcategoryCode string // Alertly subcategory
	Latitude       float64
	Longitude      float64
	ImageURL       string  // Official asset or camera image
	Address        string  // Geocoded address
	City           string
	Province       string
	PostalCode     string
	EventType      string  // Category name (e.g., "crime", "fire_incident")
}

// BotIncidentHash represents a deduplication record
type BotIncidentHash struct {
	Hash         string
	Source       string
	ExternalID   string
	CategoryCode string
	CreatedAt    time.Time
	ExpiresAt    *time.Time
}

// GeocodingCache represents a cached geocoding result
type GeocodingCache struct {
	AddressHash        string
	OriginalAddress    string
	NormalizedAddress  string
	Latitude           float64
	Longitude          float64
	Source             string
	CreatedAt          time.Time
	LastUsedAt         time.Time
}

// OfficialAssetsMap maps Alertly categories to S3 official report images
var OfficialAssetsMap = map[string]string{
	"crime":                        "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/crime.webp",
	"traffic_accident":             "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/traffic_accident.webp",
	"medical_emergency":            "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/medical_emergency.webp",
	"fire_incident":                "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/fire_incident.webp",
	"vandalism":                    "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/vandalism.webp",
	"suspicious_activity":          "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/suspicious_activity.webp",
	"infrastructure_issues":        "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/infrastructure_issue.webp",
	"extreme_weather":              "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/extreme_weather.webp",
	"community_events":             "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/community_event.webp",
	"dangerous_wildlife_sighting":  "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/dangerous_wildlife.webp",
	"positive_actions":             "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/positive_action.webp",
	"lost_pet":                     "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/lost_pet.webp",
}

// GetOfficialAsset returns the S3 URL for a category's official report image
func GetOfficialAsset(categoryCode string) string {
	if url, exists := OfficialAssetsMap[categoryCode]; exists {
		return url
	}
	// Fallback to generic incident image
	return "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/crime.webp"
}
