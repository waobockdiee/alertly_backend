package cjbot_creator

import (
	"alertly/internal/cronjobs/cjbot_creator/scrapers"
	"fmt"
	"regexp"
	"strings"
)

// CategoryMapping represents a mapping rule from source category to Alertly category
type CategoryMapping struct {
	Pattern         *regexp.Regexp
	CategoryCode    string
	SubcategoryCode string
	Priority        int // Higher priority wins if multiple matches
}

// Normalizer handles conversion of source-specific data to Alertly schema
type Normalizer struct {
	mappings map[string][]CategoryMapping // key: source name
}

// NewNormalizer creates a normalizer with all mapping rules
func NewNormalizer() *Normalizer {
	n := &Normalizer{
		mappings: make(map[string][]CategoryMapping),
	}
	n.initTPSMappings()
	n.initTFSMappings()
	n.initTTCMappings()
	n.initHydroMappings()
	n.initWeatherMappings()
	return n
}

// ============================================
// TPS (TORONTO POLICE SERVICE) MAPPINGS
// ============================================

func (n *Normalizer) initTPSMappings() {
	n.mappings["tps"] = []CategoryMapping{
		// ✅ Traffic accidents (using EXACT DB codes)
		{regexp.MustCompile(`(?i)PIACC|PERSONAL.*INJURY.*COLLISION`), "traffic_accident", "pedestrian_nvolvement", 20}, // ⚠️ Frontend typo
		{regexp.MustCompile(`(?i)PDACC|PROPERTY.*DAMAGE.*COLLISION`), "traffic_accident", "single_vehicle_accident", 18},
		{regexp.MustCompile(`(?i)FTRPD|FAIL.*REMAIN`), "traffic_accident", "hit_and_run", 19},
		{regexp.MustCompile(`(?i)IMPDR|IMPAIRED.*DRIVER`), "traffic_accident", "single_vehicle_accident", 17},
		{regexp.MustCompile(`(?i)TRAOB|TRAFFIC.*OBSTRUCT`), "traffic_accident", "single_vehicle_accident", 15},

		// ✅ Violent crime (FIXED: using valid codes)
		{regexp.MustCompile(`(?i)ASSPR|ASSAULT.*PROGRESS`), "crime", "assault", 26},
		{regexp.MustCompile(`(?i)ASSJU|ASSAULT.*JUST.*OCCURRED`), "crime", "assault", 25},
		{regexp.MustCompile(`(?i)ASS|ASSAULT`), "crime", "assault", 24},
		{regexp.MustCompile(`(?i)ROB|ROBBERY`), "crime", "robbery", 28},
		{regexp.MustCompile(`(?i)PERGU|PERSON.*GUN`), "crime", "assault", 30}, // High priority
		{regexp.MustCompile(`(?i)BREPR|BREAK.*ENTER.*PROGRESS`), "crime", "robbery", 22},
		{regexp.MustCompile(`(?i)ATTBR|ATTEMPT.*BREAK.*ENTER`), "crime", "robbery", 21},
		{regexp.MustCompile(`(?i)BREEN|BREAK.*ENTER`), "crime", "robbery", 20},

		// ✅ Theft (FIXED)
		{regexp.MustCompile(`(?i)THEJU|THEFT.*JUST.*OCCURRED`), "crime", "theft", 23},
		{regexp.MustCompile(`(?i)THE\b|THEFT`), "crime", "theft", 22},
		{regexp.MustCompile(`(?i)FRA|FRAUD`), "crime", "fraud", 24},

		// ✅ Property damage (FIXED: vandalism category exists)
		{regexp.MustCompile(`(?i)DAMJU|DAMAGE.*JUST.*OCCURRED`), "vandalism", "public_property_damage", 16},
		{regexp.MustCompile(`(?i)HAZ|HAZARD`), "infrastructure_issues", "public_utility_issues", 18},

		// ✅ Disturbances (already correct)
		{regexp.MustCompile(`(?i)DIS|DISORDERLIES`), "suspicious_activity", "unusual_behavior", 10},
		{regexp.MustCompile(`(?i)DISPU|DISPUTE`), "suspicious_activity", "unusual_behavior", 10},

		// ✅ Medical/Ambulance (FIXED: using valid medical codes)
		{regexp.MustCompile(`(?i)SEEAM|SEE.*AMBULANCE`), "medical_emergency", "other_medical_emergency", 12},
		{regexp.MustCompile(`(?i)SEEFI|SEE.*FIRE`), "fire_incident", "other_fire_incident", 12},
		{regexp.MustCompile(`(?i)FIR|FIRE`), "fire_incident", "other_fire_incident", 15},

		// ✅ Other
		{regexp.MustCompile(`(?i)UNKTR|UNKNOWN.*TROUBLE`), "suspicious_activity", "unusual_behavior", 5},
		{regexp.MustCompile(`(?i)ARR|ARREST`), "crime", "assault", 15},
		{regexp.MustCompile(`(?i)ANICO|ANIMAL.*COMP`), "dangerous_wildlife_sighting", "other_wildlife", 12},
		{regexp.MustCompile(`(?i)TAXAL|TAXI.*ALARM`), "suspicious_activity", "unusual_behavior", 8},

		// ✅ Fallback (already correct)
		{regexp.MustCompile(`(?i).*`), "suspicious_activity", "unusual_behavior", 1},
	}
}

// ============================================
// TFS (TORONTO FIRE SERVICES) MAPPINGS
// ============================================

func (n *Normalizer) initTFSMappings() {
	n.mappings["tfs"] = []CategoryMapping{
		// Fire incidents - using VALID subcategory codes from Categories.tsx
		{regexp.MustCompile(`(?i)fire.*residential|residential.*fire`), "fire_incident", "residential_fire", 15},
		{regexp.MustCompile(`(?i)fire.*highrise|highrise.*fire`), "fire_incident", "residential_fire", 16},
		{regexp.MustCompile(`(?i)fire.*commercial|commercial.*fire`), "fire_incident", "residential_fire", 15},
		{regexp.MustCompile(`(?i)vehicle.*fire|car.*fire`), "fire_incident", "vehicle_fire", 15},
		{regexp.MustCompile(`(?i)alarm.*single.*source`), "fire_incident", "other_fire_incident", 12},
		{regexp.MustCompile(`(?i)alarm|detector`), "fire_incident", "other_fire_incident", 10},
		{regexp.MustCompile(`(?i)fire|burning|smoke|flames`), "fire_incident", "residential_fire", 8},

		// Medical emergencies - using VALID subcategory codes from Categories.tsx
		{regexp.MustCompile(`(?i)^medical$`), "medical_emergency", "other_medical_emergency", 10},
		{regexp.MustCompile(`(?i)overdose|poisoning`), "medical_emergency", "overdose_poisoning", 15},
		{regexp.MustCompile(`(?i)cardiac|heart`), "medical_emergency", "cardiac_arrest", 12},
		{regexp.MustCompile(`(?i)stroke`), "medical_emergency", "stroke", 12},
		{regexp.MustCompile(`(?i)trauma|injury`), "medical_emergency", "trauma_Injury", 12},

		// Vehicle accidents (using EXACT DB codes)
		{regexp.MustCompile(`(?i)vehicle.*personal.*injury`), "traffic_accident", "single_vehicle_accident", 15},
		{regexp.MustCompile(`(?i)vehicle.*highway`), "traffic_accident", "single_vehicle_accident", 14},

		// Hazmat
		{regexp.MustCompile(`(?i)hazmat|chemical|gas.*leak|spill`), "fire_incident", "other_fire_incident", 20},

		// Rescue
		{regexp.MustCompile(`(?i)rescue|trapped|confined`), "fire_incident", "other_fire_incident", 12},

		// Fallback
		{regexp.MustCompile(`(?i).*`), "fire_incident", "other_fire_incident", 1},
	}
}

// ============================================
// TTC (TRANSIT) MAPPINGS
// ============================================

func (n *Normalizer) initTTCMappings() {
	n.mappings["ttc"] = []CategoryMapping{
		// Infrastructure issues - using VALID subcategory codes from Categories.tsx
		{regexp.MustCompile(`(?i)delay|service|suspended|disruption`), "infrastructure_issues", "public_utility_issues", 10},
		{regexp.MustCompile(`(?i)signal|mechanical|technical`), "infrastructure_issues", "public_utility_issues", 8},
		{regexp.MustCompile(`(?i)power|electrical|outage`), "infrastructure_issues", "public_utility_issues", 12},

		// Medical emergencies on transit - using VALID subcategory codes from Categories.tsx
		{regexp.MustCompile(`(?i)medical|emergency.*ttc|passenger.*injury`), "medical_emergency", "trauma_Injury", 15},

		// Security incidents - using VALID subcategory codes from Categories.tsx
		{regexp.MustCompile(`(?i)security|police|investigation`), "suspicious_activity", "suspicious_person", 10},

		// Fallback
		{regexp.MustCompile(`(?i).*`), "infrastructure_issues", "public_utility_issues", 1},
	}
}

// ============================================
// TORONTO HYDRO (POWER OUTAGES) MAPPINGS
// ============================================

func (n *Normalizer) initHydroMappings() {
	n.mappings["hydro"] = []CategoryMapping{
		// All hydro incidents map to infrastructure/utility
		{regexp.MustCompile(`(?i)outage|power.*out|no.*power`), "infrastructure_issues", "public_utility_issues", 10},
		{regexp.MustCompile(`(?i)planned|maintenance|scheduled`), "infrastructure_issues", "public_utility_issues", 8},
		{regexp.MustCompile(`(?i)unplanned|emergency|fault`), "infrastructure_issues", "public_utility_issues", 12},

		// Fallback
		{regexp.MustCompile(`(?i).*`), "infrastructure_issues", "public_utility_issues", 1},
	}
}

// ============================================
// ENVIRONMENT CANADA (WEATHER) MAPPINGS
// ============================================

func (n *Normalizer) initWeatherMappings() {
	n.mappings["weather"] = []CategoryMapping{
		// Extreme weather categories
		{regexp.MustCompile(`(?i)tornado|funnel.*cloud`), "extreme_weather", "high_winds_tornado", 30},
		{regexp.MustCompile(`(?i)snow.*storm|blizzard|winter.*storm`), "extreme_weather", "snow_storm", 20},
		{regexp.MustCompile(`(?i)thunderstorm|lightning|severe.*storm`), "extreme_weather", "hail_severe_storm", 15},
		{regexp.MustCompile(`(?i)flood|flooding|heavy.*rain`), "extreme_weather", "heavy_rain_flooding", 18},
		{regexp.MustCompile(`(?i)heat|extreme.*temp|hot`), "extreme_weather", "extreme_heat", 12},
		{regexp.MustCompile(`(?i)cold|freeze|frost|wind.*chill`), "extreme_weather", "icy_roads", 12},
		{regexp.MustCompile(`(?i)wind|gale|storm.*wind`), "extreme_weather", "high_winds_tornado", 10},
		{regexp.MustCompile(`(?i)hail`), "extreme_weather", "hail_severe_storm", 15},

		// Fallback
		{regexp.MustCompile(`(?i)weather|advisory|warning`), "extreme_weather", "hail_severe_storm", 1},
	}
}

// ============================================
// NORMALIZATION METHOD
// ============================================

// Normalize converts a ScrapedIncident to NormalizedIncident
func (n *Normalizer) Normalize(scraped scrapers.ScrapedIncident) (*NormalizedIncident, error) {
	mappings, exists := n.mappings[scraped.Source]
	if !exists {
		return nil, fmt.Errorf("no mappings for source: %s", scraped.Source)
	}

	// Find best matching mapping
	var bestMatch *CategoryMapping
	highestPriority := -1

	for i := range mappings {
		m := &mappings[i]
		if m.Pattern.MatchString(scraped.RawCategory) || m.Pattern.MatchString(scraped.RawTitle) {
			if m.Priority > highestPriority {
				bestMatch = m
				highestPriority = m.Priority
			}
		}
	}

	if bestMatch == nil {
		return nil, fmt.Errorf("no category match found for: %s", scraped.RawCategory)
	}

	// Build normalized incident
	normalized := &NormalizedIncident{
		Title:           n.buildTitle(scraped, bestMatch),
		Description:     n.buildDescription(scraped),
		CategoryCode:    bestMatch.CategoryCode,
		SubcategoryCode: bestMatch.SubcategoryCode,
		EventType:       bestMatch.CategoryCode, // event_type = category name
	}

	// Determine image URL (camera URL or official asset)
	if scraped.CameraURL != "" && strings.HasSuffix(scraped.CameraURL, ".jpg") {
		normalized.ImageURL = scraped.CameraURL
	} else {
		normalized.ImageURL = GetOfficialAsset(bestMatch.CategoryCode)
	}

	// Handle coordinates
	if scraped.Latitude != nil && scraped.Longitude != nil {
		normalized.Latitude = *scraped.Latitude
		normalized.Longitude = *scraped.Longitude
	} else if len(scraped.Polygon) > 0 {
		// Calculate centroid for polygon-based incidents
		centroid := CalculateCentroid(scraped.Polygon)
		normalized.Latitude = centroid.Lat
		normalized.Longitude = centroid.Lng
	}
	// Note: If no coordinates available, they will be geocoded later from address
	// normalized.Latitude and normalized.Longitude will be 0 (will be checked in service)

	return normalized, nil
}

// buildTitle creates a user-friendly title
func (n *Normalizer) buildTitle(scraped scrapers.ScrapedIncident, mapping *CategoryMapping) string {
	if scraped.RawTitle != "" {
		return scraped.RawTitle
	}

	// Generate title from category
	switch scraped.Source {
	case "tps":
		return fmt.Sprintf("Police Call: %s", scraped.RawCategory)
	case "tfs":
		return fmt.Sprintf("Fire Service Call: %s", scraped.RawCategory)
	case "ttc":
		return fmt.Sprintf("Transit Alert: %s", scraped.RawCategory)
	case "hydro":
		return "Power Outage"
	case "weather":
		return fmt.Sprintf("Weather Alert: %s", scraped.RawCategory)
	default:
		return "Official Report"
	}
}

// buildDescription creates detailed narrative description
func (n *Normalizer) buildDescription(scraped scrapers.ScrapedIncident) string {
	var narrative string

	// Build narrative based on source
	switch scraped.Source {
	case "tps":
		narrative = n.buildTPSNarrative(scraped)
	case "tfs":
		narrative = n.buildTFSNarrative(scraped)
	case "ttc":
		narrative = n.buildTTCNarrative(scraped)
	case "weather":
		narrative = n.buildWeatherNarrative(scraped)
	default:
		narrative = scraped.RawDescription
	}

	// Add source attribution (full name)
	sourceName := n.getSourceFullName(scraped.Source)
	narrative += fmt.Sprintf("\n\nSource: %s", sourceName)

	return narrative
}

// getSourceFullName returns the full official name of the data source
func (n *Normalizer) getSourceFullName(source string) string {
	switch source {
	case "tps":
		return "Toronto Police Service"
	case "tfs":
		return "Toronto Fire Services"
	case "ttc":
		return "TTC"
	case "hydro":
		return "Toronto Hydro"
	case "weather":
		return "Environment Canada"
	default:
		return strings.ToUpper(source)
	}
}

// buildTPSNarrative creates a narrative description for police calls
func (n *Normalizer) buildTPSNarrative(scraped scrapers.ScrapedIncident) string {
	callType := scraped.RawTitle
	location := scraped.Address
	if location == "" || location == "Toronto, ON" {
		location = "an undisclosed location in Toronto"
	}

	// Create natural language narrative based on call type
	var narrative string
	code := strings.ToUpper(scraped.RawCategory)

	switch {
	case strings.Contains(code, "PIACC"):
		narrative = fmt.Sprintf("Officers responded to a personal injury collision at %s.", location)
	case strings.Contains(code, "PDACC"):
		narrative = fmt.Sprintf("A property damage collision was reported at %s.", location)
	case strings.Contains(code, "ASSJU"):
		narrative = fmt.Sprintf("Police attended an assault that just occurred at %s.", location)
	case strings.Contains(code, "ROB"):
		narrative = fmt.Sprintf("A robbery was reported at %s. Officers are investigating.", location)
	case strings.Contains(code, "BREPR"):
		narrative = fmt.Sprintf("A break and enter in progress was reported at %s. Police units responded to the scene.", location)
	case strings.Contains(code, "THE"):
		narrative = fmt.Sprintf("A theft incident was reported at %s.", location)
	case strings.Contains(code, "DIS"):
		narrative = fmt.Sprintf("Officers responded to a disturbance at %s.", location)
	case strings.Contains(code, "MISEL"):
		narrative = fmt.Sprintf("Police are searching for a missing elderly person last seen near %s.", location)
	case strings.Contains(code, "SEEAM") || strings.Contains(code, "ASSAM"):
		narrative = fmt.Sprintf("Emergency medical services were requested at %s.", location)
	default:
		narrative = fmt.Sprintf("%s reported at %s.", callType, location)
	}

	return narrative
}

// buildTFSNarrative creates a narrative description for fire service calls
func (n *Normalizer) buildTFSNarrative(scraped scrapers.ScrapedIncident) string {
	eventType := scraped.RawTitle
	location := scraped.Address
	if location == "" || location == "Toronto, ON" {
		location = "a location in Toronto"
	}

	var narrative string
	category := strings.ToLower(eventType)

	// Create detailed human-readable narratives based on incident type
	switch {
	case strings.Contains(category, "structure fire"):
		narrative = fmt.Sprintf("Toronto Fire Services responded to a structure fire at %s. Multiple units have been dispatched to the scene.", location)
	case strings.Contains(category, "residential fire") || strings.Contains(category, "fire residential"):
		narrative = fmt.Sprintf("Fire crews are responding to a residential fire at %s. Emergency services are on scene.", location)
	case strings.Contains(category, "highrise fire") || strings.Contains(category, "fire highrise"):
		narrative = fmt.Sprintf("A fire has been reported in a highrise building at %s. Multiple fire companies are responding.", location)
	case strings.Contains(category, "commercial fire") || strings.Contains(category, "fire commercial"):
		narrative = fmt.Sprintf("Fire units are responding to a commercial building fire at %s.", location)
	case strings.Contains(category, "vehicle fire") || strings.Contains(category, "car fire"):
		narrative = fmt.Sprintf("A vehicle fire has been reported at %s. Fire crews are responding to extinguish the blaze.", location)
	case strings.Contains(category, "medical"):
		narrative = fmt.Sprintf("Toronto Fire paramedics have been dispatched to a medical emergency at %s.", location)
	case strings.Contains(category, "vehicle") && strings.Contains(category, "personal injury"):
		narrative = fmt.Sprintf("Fire and medical units are responding to a vehicle collision with injuries at %s.", location)
	case strings.Contains(category, "vehicle") && strings.Contains(category, "highway"):
		narrative = fmt.Sprintf("Emergency crews are responding to a vehicle incident on the highway near %s.", location)
	case strings.Contains(category, "alarm"):
		if strings.Contains(category, "single") {
			narrative = fmt.Sprintf("Fire crews are investigating a single-source alarm activation at %s.", location)
		} else {
			narrative = fmt.Sprintf("Fire services are responding to an alarm activation at %s.", location)
		}
	case strings.Contains(category, "hazmat") || strings.Contains(category, "chemical") || strings.Contains(category, "gas leak"):
		narrative = fmt.Sprintf("Hazardous materials teams are responding to a reported %s at %s.", strings.ToLower(eventType), location)
	case strings.Contains(category, "rescue"):
		narrative = fmt.Sprintf("Specialized rescue teams have been dispatched to %s for a technical rescue operation.", location)
	default:
		narrative = fmt.Sprintf("Toronto Fire Services are responding to a %s at %s.", strings.ToLower(eventType), location)
	}

	return narrative
}

// buildTTCNarrative creates a narrative description for transit alerts
func (n *Normalizer) buildTTCNarrative(scraped scrapers.ScrapedIncident) string {
	return scraped.RawDescription
}

// buildWeatherNarrative creates a narrative description for weather alerts
func (n *Normalizer) buildWeatherNarrative(scraped scrapers.ScrapedIncident) string {
	return scraped.RawDescription
}

// CalculateCentroid computes the geographic center of a polygon
func CalculateCentroid(points []scrapers.Point) scrapers.Point {
	if len(points) == 0 {
		return scrapers.Point{Lat: 43.6532, Lng: -79.3832} // Toronto center fallback
	}

	var sumLat, sumLng float64
	for _, p := range points {
		sumLat += p.Lat
		sumLng += p.Lng
	}

	return scrapers.Point{
		Lat: sumLat / float64(len(points)),
		Lng: sumLng / float64(len(points)),
	}
}
