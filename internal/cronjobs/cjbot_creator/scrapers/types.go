package scrapers

import "time"

// ScrapedIncident represents raw data from any external source before normalization
type ScrapedIncident struct {
	Source         string    // "tps", "tfs", "ttc", "hydro", "weather"
	ExternalID     string    // Unique ID from source system
	RawTitle       string    // Original title/description
	RawDescription string    // Additional details
	RawCategory    string    // Source's category (e.g., "ASSAULT", "MEDICAL CALL")
	RawSubcategory string    // Source's subcategory (optional)
	Address        string    // Text address (for TPS/TFS)
	Latitude       *float64  // Coordinates (if provided by source)
	Longitude      *float64  // Coordinates (if provided by source)
	Polygon        []Point   // For hydro outages (area affected)
	Timestamp      time.Time // When incident occurred
	ETR            *time.Time // Estimated time of restoration (for outages)
	Status         string    // "active", "resolved", etc.
	CameraURL      string    // Traffic camera URL if available
}

// Point represents a geographic coordinate
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
