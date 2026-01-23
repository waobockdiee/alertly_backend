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

// OfficialReportImageURL is the generic image used for all bot-created incidents
// This image indicates the report comes from official government data sources
const OfficialReportImageURL = "https://images.alertly.ca/official_picture.png"

// GetOfficialAsset returns the official report image URL for bot-created incidents
func GetOfficialAsset(categoryCode string) string {
	return OfficialReportImageURL
}
