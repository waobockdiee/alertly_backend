package ai

import (
	"database/sql"
	"fmt"
	"time"
)

// SmartAlertSystem provides intelligent, personalized alerts
type SmartAlertSystem struct {
	db *sql.DB
}

// SmartAlert represents an intelligent alert
type SmartAlert struct {
	AlertID        int64     `json:"alert_id"`
	UserID         int64     `json:"user_id"`
	Type           string    `json:"type"` // "prediction", "pattern", "safety", "community"
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	Priority       string    `json:"priority"` // "low", "medium", "high", "urgent"
	Category       string    `json:"category"`
	Location       Location  `json:"location"`
	Prediction     *IncidentPrediction `json:"prediction,omitempty"`
	ActionRequired bool      `json:"action_required"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    int     `json:"radius"`
	Name      string  `json:"name"`
}

// NewSmartAlertSystem creates a new smart alert system
func NewSmartAlertSystem(db *sql.DB) *SmartAlertSystem {
	return &SmartAlertSystem{db: db}
}

// GeneratePersonalizedAlerts creates alerts based on user behavior and location
func (sas *SmartAlertSystem) GeneratePersonalizedAlerts(userID int64) ([]SmartAlert, error) {
	var alerts []SmartAlert
	
	// Get user preferences and behavior
	userProfile, err := sas.getUserProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	
	// Get user's favorite locations
	favoriteLocations, err := sas.getFavoriteLocations(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite locations: %w", err)
	}
	
	// Generate location-based alerts
	for _, location := range favoriteLocations {
		locationAlerts, err := sas.generateLocationAlerts(userID, location, userProfile)
		if err != nil {
			continue // Skip this location if there's an error
		}
		alerts = append(alerts, locationAlerts...)
	}
	
	// Generate pattern-based alerts
	patternAlerts, err := sas.generatePatternAlerts(userID, userProfile)
	if err == nil {
		alerts = append(alerts, patternAlerts...)
	}
	
	// Generate safety alerts
	safetyAlerts, err := sas.generateSafetyAlerts(userID, userProfile)
	if err == nil {
		alerts = append(alerts, safetyAlerts...)
	}
	
	return alerts, nil
}

// generateLocationAlerts creates alerts for specific locations
func (sas *SmartAlertSystem) generateLocationAlerts(userID int64, location Location, profile UserProfile) ([]SmartAlert, error) {
	var alerts []SmartAlert
	
	// Analyze risk for this location
	predictiveAnalytics := NewPredictiveAnalytics(sas.db)
	prediction, err := predictiveAnalytics.AnalyzeAreaRisk(location.Latitude, location.Longitude, location.Radius)
	if err != nil {
		return alerts, err
	}
	
	// Generate alerts based on risk level
	switch prediction.RiskLevel {
	case "critical":
		alerts = append(alerts, SmartAlert{
			UserID:   userID,
			Type:     "prediction",
			Title:    "🚨 Critical Risk Alert",
			Message:  fmt.Sprintf("High risk detected in %s. Consider avoiding the area or taking extra precautions.", location.Name),
			Priority: "urgent",
			Category: "safety",
			Location: location,
			Prediction: prediction,
			ActionRequired: true,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now(),
		})
	case "high":
		alerts = append(alerts, SmartAlert{
			UserID:   userID,
			Type:     "prediction",
			Title:    "⚠️ High Risk Warning",
			Message:  fmt.Sprintf("Elevated risk detected in %s. Stay alert and report any suspicious activity.", location.Name),
			Priority: "high",
			Category: "safety",
			Location: location,
			Prediction: prediction,
			ActionRequired: false,
			ExpiresAt: time.Now().Add(12 * time.Hour),
			CreatedAt: time.Now(),
		})
	}
	
	// Generate category-specific alerts based on user preferences
	for _, predictedIncident := range prediction.PredictedIncidents {
		if sas.isCategoryRelevant(predictedIncident.Category, profile) {
			alerts = append(alerts, SmartAlert{
				UserID:   userID,
				Type:     "prediction",
				Title:    fmt.Sprintf("🔮 %s Prediction", sas.getCategoryDisplayName(predictedIncident.Category)),
				Message:  fmt.Sprintf("Potential %s incident predicted in %s within %s.", 
					sas.getCategoryDisplayName(predictedIncident.Category), 
					location.Name, 
					predictedIncident.TimeWindow),
				Priority: "medium",
				Category: predictedIncident.Category,
				Location: location,
				ActionRequired: false,
				ExpiresAt: time.Now().Add(6 * time.Hour),
				CreatedAt: time.Now(),
			})
		}
	}
	
	return alerts, nil
}

// generatePatternAlerts creates alerts based on user behavior patterns
func (sas *SmartAlertSystem) generatePatternAlerts(userID int64, profile UserProfile) ([]SmartAlert, error) {
	var alerts []SmartAlert
	
	// Analyze user's reporting patterns
	patterns, err := sas.analyzeUserPatterns(userID)
	if err != nil {
		return alerts, err
	}
	
	// Generate alerts based on patterns
	if patterns.ReportFrequency > 5 && patterns.LastReportDate.Before(time.Now().AddDate(0, 0, -7)) {
		alerts = append(alerts, SmartAlert{
			UserID:   userID,
			Type:     "pattern",
			Title:    "📊 Community Update",
			Message:  "You've been quiet lately! Your community could benefit from your active reporting. Consider checking your usual areas.",
			Priority: "low",
			Category: "community",
			ActionRequired: false,
			ExpiresAt: time.Now().AddDate(0, 0, 3),
			CreatedAt: time.Now(),
		})
	}
	
	// Alert about new badges or achievements
	if patterns.PotentialBadges > 0 {
		alerts = append(alerts, SmartAlert{
			UserID:   userID,
			Type:     "pattern",
			Title:    "🏆 Achievement Alert",
			Message:  fmt.Sprintf("You're close to earning %d new badge(s)! Keep up the great work.", patterns.PotentialBadges),
			Priority: "low",
			Category: "community",
			ActionRequired: false,
			ExpiresAt: time.Now().AddDate(0, 0, 7),
			CreatedAt: time.Now(),
		})
	}
	
	return alerts, nil
}

// generateSafetyAlerts creates safety-focused alerts
func (sas *SmartAlertSystem) generateSafetyAlerts(userID int64, profile UserProfile) ([]SmartAlert, error) {
	var alerts []SmartAlert
	
	// Check for safety patterns
	safetyScore, err := sas.calculateSafetyScore(userID)
	if err != nil {
		return alerts, err
	}
	
	if safetyScore < 0.3 {
		alerts = append(alerts, SmartAlert{
			UserID:   userID,
			Type:     "safety",
			Title:    "🛡️ Safety Reminder",
			Message:  "Your safety score is lower than usual. Consider reviewing recent incidents in your area and staying extra vigilant.",
			Priority: "medium",
			Category: "safety",
			ActionRequired: false,
			ExpiresAt: time.Now().AddDate(0, 0, 1),
			CreatedAt: time.Now(),
		})
	}
	
	return alerts, nil
}

// Helper methods
type UserProfile struct {
	UserID           int64     `json:"user_id"`
	Credibility      float64   `json:"credibility"`
	ReportFrequency  int       `json:"report_frequency"`
	LastReportDate   time.Time `json:"last_report_date"`
	PreferredCategories []string `json:"preferred_categories"`
	SafetyScore      float64   `json:"safety_score"`
}

type UserPatterns struct {
	ReportFrequency   int       `json:"report_frequency"`
	LastReportDate    time.Time `json:"last_report_date"`
	PotentialBadges   int       `json:"potential_badges"`
	ActiveAreas       []string  `json:"active_areas"`
}

func (sas *SmartAlertSystem) getUserProfile(userID int64) (UserProfile, error) {
	query := `
		SELECT 
			account_id,
			credibility,
			counter_total_incidents_created,
			last_incident_date
		FROM account
		WHERE account_id = $1
	`

	var profile UserProfile
	err := sas.db.QueryRow(query, userID).Scan(
		&profile.UserID,
		&profile.Credibility,
		&profile.ReportFrequency,
		&profile.LastReportDate,
	)
	
	if err != nil {
		return profile, err
	}
	
	// Get preferred categories
	profile.PreferredCategories, _ = sas.getPreferredCategories(userID)
	
	return profile, nil
}

func (sas *SmartAlertSystem) getFavoriteLocations(userID int64) ([]Location, error) {
	query := `
		SELECT
			title,
			latitude,
			longitude,
			radius
		FROM account_favorite_locations
		WHERE account_id = $1
	`

	rows, err := sas.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var locations []Location
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.Name, &location.Latitude, &location.Longitude, &location.Radius)
		if err != nil {
			continue
		}
		locations = append(locations, location)
	}
	
	return locations, nil
}

func (sas *SmartAlertSystem) analyzeUserPatterns(userID int64) (UserPatterns, error) {
	query := `
		SELECT
			COUNT(*) as report_count,
			MAX(created_at) as last_report
		FROM incident_reports
		WHERE account_id = $1
		AND created_at >= NOW() - INTERVAL '30 days'
	`

	var patterns UserPatterns
	err := sas.db.QueryRow(query, userID).Scan(&patterns.ReportFrequency, &patterns.LastReportDate)
	if err != nil {
		return patterns, err
	}
	
	// Calculate potential badges (simplified)
	patterns.PotentialBadges = sas.calculatePotentialBadges(userID)
	
	return patterns, nil
}

func (sas *SmartAlertSystem) calculateSafetyScore(userID int64) (float64, error) {
	// Implement safety score calculation based on user's incident history
	// and surrounding area safety
	return 0.7, nil // Placeholder
}

func (sas *SmartAlertSystem) getPreferredCategories(userID int64) ([]string, error) {
	query := `
		SELECT DISTINCT category_code
		FROM incident_reports
		WHERE account_id = $1
		ORDER BY COUNT(*) DESC
		LIMIT 5
	`
	
	rows, err := sas.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err == nil {
			categories = append(categories, category)
		}
	}
	
	return categories, nil
}

func (sas *SmartAlertSystem) isCategoryRelevant(category string, profile UserProfile) bool {
	for _, preferred := range profile.PreferredCategories {
		if preferred == category {
			return true
		}
	}
	return false
}

func (sas *SmartAlertSystem) getCategoryDisplayName(category string) string {
	displayNames := map[string]string{
		"crime": "Crime",
		"traffic_accident": "Traffic Accident",
		"medical_emergency": "Medical Emergency",
		"fire_incident": "Fire Incident",
		"vandalism": "Vandalism",
		"suspicious_activity": "Suspicious Activity",
		"infrastructure_issues": "Infrastructure Issue",
		"extreme_weather": "Extreme Weather",
		"community_events": "Community Event",
		"dangerous_wildlife_sighting": "Wildlife Sighting",
		"positive_actions": "Positive Action",
		"lost_pet": "Lost Pet",
	}
	
	if name, exists := displayNames[category]; exists {
		return name
	}
	return category
}

func (sas *SmartAlertSystem) calculatePotentialBadges(userID int64) int {
	// Implement badge calculation logic
	return 2 // Placeholder
}
