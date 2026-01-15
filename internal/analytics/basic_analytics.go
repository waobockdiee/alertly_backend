package analytics

import (
	"database/sql"
	"fmt"
	"time"
)

// BasicAnalytics provides simple analytics that work with limited data
type BasicAnalytics struct {
	db *sql.DB
}

// AnalyticsSummary represents basic analytics data
type AnalyticsSummary struct {
	TotalIncidents     int              `json:"total_incidents"`
	TotalUsers         int              `json:"total_users"`
	ActiveUsers        int              `json:"active_users"`
	AverageCredibility float64          `json:"average_credibility"`
	TopCategories      []CategoryStats  `json:"top_categories"`
	RecentActivity     []RecentActivity `json:"recent_activity"`
	DataQuality        DataQuality      `json:"data_quality"`
}

type CategoryStats struct {
	Category   string  `json:"category"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type RecentActivity struct {
	Date          time.Time `json:"date"`
	IncidentCount int       `json:"incident_count"`
	UserCount     int       `json:"user_count"`
}

type DataQuality struct {
	HasMinimumData   bool     `json:"has_minimum_data"`
	DataCompleteness float64  `json:"data_completeness"`
	Recommendations  []string `json:"recommendations"`
}

// NewBasicAnalytics creates a new basic analytics service
func NewBasicAnalytics(db *sql.DB) *BasicAnalytics {
	return &BasicAnalytics{db: db}
}

// GetAnalyticsSummary returns basic analytics summary
func (ba *BasicAnalytics) GetAnalyticsSummary() (*AnalyticsSummary, error) {
	summary := &AnalyticsSummary{}

	// Get total incidents
	totalIncidents, err := ba.getTotalIncidents()
	if err != nil {
		return nil, err
	}
	summary.TotalIncidents = totalIncidents

	// Get total users
	totalUsers, err := ba.getTotalUsers()
	if err != nil {
		return nil, err
	}
	summary.TotalUsers = totalUsers

	// Get active users (last 30 days)
	activeUsers, err := ba.getActiveUsers()
	if err != nil {
		return nil, err
	}
	summary.ActiveUsers = activeUsers

	// Get average credibility
	avgCredibility, err := ba.getAverageCredibility()
	if err != nil {
		return nil, err
	}
	summary.AverageCredibility = avgCredibility

	// Get top categories
	topCategories, err := ba.getTopCategories()
	if err != nil {
		return nil, err
	}
	summary.TopCategories = topCategories

	// Get recent activity
	recentActivity, err := ba.getRecentActivity()
	if err != nil {
		return nil, err
	}
	summary.RecentActivity = recentActivity

	// Assess data quality
	dataQuality, err := ba.assessDataQuality()
	if err != nil {
		return nil, err
	}
	summary.DataQuality = dataQuality

	return summary, nil
}

// GetLocationAnalytics returns analytics for a specific location
func (ba *BasicAnalytics) GetLocationAnalytics(latitude, longitude float64, radiusMeters int) (*LocationAnalytics, error) {
	analytics := &LocationAnalytics{
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radiusMeters,
	}

	// Get incidents in this area
	incidents, err := ba.getIncidentsInArea(latitude, longitude, radiusMeters)
	if err != nil {
		return nil, err
	}
	analytics.TotalIncidents = incidents

	// Get active users in this area
	activeUsers, err := ba.getActiveUsersInArea(latitude, longitude, radiusMeters)
	if err != nil {
		return nil, err
	}
	analytics.ActiveUsers = activeUsers

	// Get recent activity in this area
	recentActivity, err := ba.getRecentActivityInArea(latitude, longitude, radiusMeters)
	if err != nil {
		return nil, err
	}
	analytics.RecentActivity = recentActivity

	// Get top categories in this area
	topCategories, err := ba.getTopCategoriesInArea(latitude, longitude, radiusMeters)
	if err != nil {
		return nil, err
	}
	analytics.TopCategories = topCategories

	// Assess data quality for this location
	dataQuality, err := ba.assessLocationDataQuality(incidents, activeUsers)
	if err != nil {
		return nil, err
	}
	analytics.DataQuality = dataQuality

	return analytics, nil
}

// LocationAnalytics represents analytics for a specific location
type LocationAnalytics struct {
	Latitude       float64             `json:"latitude"`
	Longitude      float64             `json:"longitude"`
	Radius         int                 `json:"radius"`
	TotalIncidents int                 `json:"total_incidents"`
	ActiveUsers    int                 `json:"active_users"`
	RecentActivity []RecentActivity    `json:"recent_activity"`
	TopCategories  []CategoryStats     `json:"top_categories"`
	DataQuality    LocationDataQuality `json:"data_quality"`
}

type LocationDataQuality struct {
	HasMinimumData   bool     `json:"has_minimum_data"`
	DataCompleteness float64  `json:"data_completeness"`
	Recommendations  []string `json:"recommendations"`
	RiskLevel        string   `json:"risk_level"` // low, medium, high, critical
}

// getTotalIncidents returns total number of incidents
func (ba *BasicAnalytics) getTotalIncidents() (int, error) {
	query := `SELECT COUNT(*) FROM incident_reports WHERE is_active = true`
	var count int
	err := ba.db.QueryRow(query).Scan(&count)
	return count, err
}

// getTotalUsers returns total number of users
func (ba *BasicAnalytics) getTotalUsers() (int, error) {
	query := `SELECT COUNT(*) FROM account WHERE status = 'active'`
	var count int
	err := ba.db.QueryRow(query).Scan(&count)
	return count, err
}

// getActiveUsers returns users active in last 30 days
func (ba *BasicAnalytics) getActiveUsers() (int, error) {
	query := `
		SELECT COUNT(DISTINCT account_id)
		FROM incident_reports
		WHERE created_at >= NOW() - INTERVAL '30 days'
	`
	var count int
	err := ba.db.QueryRow(query).Scan(&count)
	return count, err
}

// getAverageCredibility returns average credibility score
func (ba *BasicAnalytics) getAverageCredibility() (float64, error) {
	query := `SELECT AVG(credibility) FROM account WHERE status = 'active' AND credibility IS NOT NULL`
	var avg sql.NullFloat64
	err := ba.db.QueryRow(query).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if avg.Valid {
		return avg.Float64, nil
	}
	return 0, nil
}

// getTopCategories returns top incident categories
func (ba *BasicAnalytics) getTopCategories() ([]CategoryStats, error) {
	query := `
		SELECT 
			category_code,
			COUNT(*) as count,
			(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM incident_reports WHERE is_active = true)) as percentage
		FROM incident_reports 
		WHERE is_active = true 
		GROUP BY category_code 
		ORDER BY count DESC 
		LIMIT 5
	`

	rows, err := ba.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryStats
	for rows.Next() {
		var cat CategoryStats
		err := rows.Scan(&cat.Category, &cat.Count, &cat.Percentage)
		if err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

// getRecentActivity returns recent activity (last 7 days)
func (ba *BasicAnalytics) getRecentActivity() ([]RecentActivity, error) {
	query := `
		SELECT
			DATE(created_at) as date,
			COUNT(*) as incident_count,
			COUNT(DISTINCT account_id) as user_count
		FROM incident_reports
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`

	rows, err := ba.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []RecentActivity
	for rows.Next() {
		var activity RecentActivity
		err := rows.Scan(&activity.Date, &activity.IncidentCount, &activity.UserCount)
		if err != nil {
			continue
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// assessDataQuality assesses if we have enough data for AI
func (ba *BasicAnalytics) assessDataQuality() (DataQuality, error) {
	quality := DataQuality{}

	// Check minimum data requirements
	totalIncidents, _ := ba.getTotalIncidents()
	totalUsers, _ := ba.getTotalUsers()
	activeUsers, _ := ba.getActiveUsers()

	// Calculate data completeness
	completeness := 0.0
	if totalIncidents >= 100 {
		completeness += 25
	}
	if totalUsers >= 50 {
		completeness += 25
	}
	if activeUsers >= 20 {
		completeness += 25
	}
	if totalIncidents >= 500 {
		completeness += 25
	}

	quality.DataCompleteness = completeness
	quality.HasMinimumData = completeness >= 50

	// Generate recommendations
	if totalIncidents < 100 {
		quality.Recommendations = append(quality.Recommendations,
			"Need at least 100 incidents for basic AI predictions")
	}
	if totalUsers < 50 {
		quality.Recommendations = append(quality.Recommendations,
			"Need at least 50 users for user behavior analysis")
	}
	if activeUsers < 20 {
		quality.Recommendations = append(quality.Recommendations,
			"Need at least 20 active users for pattern recognition")
	}

	return quality, nil
}

// GetSimplePredictions returns basic predictions based on current data
func (ba *BasicAnalytics) GetSimplePredictions() (map[string]interface{}, error) {
	predictions := make(map[string]interface{})

	// Only provide predictions if we have minimum data
	quality, err := ba.assessDataQuality()
	if err != nil {
		return predictions, err
	}

	if !quality.HasMinimumData {
		predictions["status"] = "insufficient_data"
		predictions["message"] = "Need more data for predictions"
		return predictions, nil
	}

	// Get basic patterns
	timePatterns, err := ba.getTimePatterns()
	if err == nil {
		predictions["time_patterns"] = timePatterns
	}

	hotspots, err := ba.getHotspots()
	if err == nil {
		predictions["hotspots"] = hotspots
	}

	predictions["status"] = "basic_predictions"
	return predictions, nil
}

// getTimePatterns returns basic time-based patterns
func (ba *BasicAnalytics) getTimePatterns() (map[string]interface{}, error) {
	query := `
		SELECT
			EXTRACT(HOUR FROM created_at)::int as hour,
			COUNT(*) as count
		FROM incident_reports
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY EXTRACT(HOUR FROM created_at)
		ORDER BY count DESC
		LIMIT 3
	`

	rows, err := ba.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patterns := make(map[string]interface{})
	var peakHours []int

	for rows.Next() {
		var hour, count int
		if err := rows.Scan(&hour, &count); err == nil {
			peakHours = append(peakHours, hour)
		}
	}

	patterns["peak_hours"] = peakHours
	patterns["message"] = fmt.Sprintf("Peak activity hours: %v", peakHours)

	return patterns, nil
}

// getHotspots returns basic geographic hotspots
func (ba *BasicAnalytics) getHotspots() (map[string]interface{}, error) {
	query := `
		SELECT
			ROUND(CAST(center_latitude AS numeric), 2) as lat,
			ROUND(CAST(center_longitude AS numeric), 2) as lng,
			COUNT(*) as count
		FROM incident_clusters
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY ROUND(CAST(center_latitude AS numeric), 2), ROUND(CAST(center_longitude AS numeric), 2)
		HAVING COUNT(*) > 1
		ORDER BY count DESC
		LIMIT 3
	`

	rows, err := ba.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hotspots := make(map[string]interface{})
	var locations []map[string]interface{}

	for rows.Next() {
		var lat, lng float64
		var count int
		if err := rows.Scan(&lat, &lng, &count); err == nil {
			locations = append(locations, map[string]interface{}{
				"latitude":  lat,
				"longitude": lng,
				"count":     count,
			})
		}
	}

	hotspots["locations"] = locations
	hotspots["message"] = fmt.Sprintf("Found %d hotspot locations", len(locations))

	return hotspots, nil
}

// getIncidentsInArea returns incidents within a specific radius
func (ba *BasicAnalytics) getIncidentsInArea(lat, lon float64, radiusMeters int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM incident_reports
		WHERE is_active = true
		AND ST_DistanceSphere(
			ST_MakePoint(longitude, latitude),
			ST_MakePoint($1, $2)
		) <= $3
	`
	var count int
	err := ba.db.QueryRow(query, lon, lat, radiusMeters).Scan(&count)
	return count, err
}

// getActiveUsersInArea returns active users within a specific radius
func (ba *BasicAnalytics) getActiveUsersInArea(lat, lon float64, radiusMeters int) (int, error) {
	query := `
		SELECT COUNT(DISTINCT ir.account_id)
		FROM incident_reports ir
		WHERE ir.created_at >= NOW() - INTERVAL '30 days'
		AND ST_DistanceSphere(
			ST_MakePoint(ir.longitude, ir.latitude),
			ST_MakePoint($1, $2)
		) <= $3
	`
	var count int
	err := ba.db.QueryRow(query, lon, lat, radiusMeters).Scan(&count)
	return count, err
}

// getRecentActivityInArea returns recent activity within a specific radius
func (ba *BasicAnalytics) getRecentActivityInArea(lat, lon float64, radiusMeters int) ([]RecentActivity, error) {
	query := `
		SELECT
			DATE(ir.created_at) as date,
			COUNT(*) as incident_count,
			COUNT(DISTINCT ir.account_id) as user_count
		FROM incident_reports ir
		WHERE ir.created_at >= NOW() - INTERVAL '7 days'
		AND ST_DistanceSphere(
			ST_MakePoint(ir.longitude, ir.latitude),
			ST_MakePoint($1, $2)
		) <= $3
		GROUP BY DATE(ir.created_at)
		ORDER BY date DESC
	`

	rows, err := ba.db.Query(query, lon, lat, radiusMeters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []RecentActivity
	for rows.Next() {
		var activity RecentActivity
		err := rows.Scan(&activity.Date, &activity.IncidentCount, &activity.UserCount)
		if err != nil {
			continue
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// getTopCategoriesInArea returns top categories within a specific radius
func (ba *BasicAnalytics) getTopCategoriesInArea(lat, lon float64, radiusMeters int) ([]CategoryStats, error) {
	query := `
		SELECT
			ir.category_code,
			COUNT(*) as count,
			(COUNT(*) * 100.0 / (
				SELECT COUNT(*)
				FROM incident_reports ir2
				WHERE ir2.is_active = true
				AND ST_DistanceSphere(
					ST_MakePoint(ir2.longitude, ir2.latitude),
					ST_MakePoint($1, $2)
				) <= $3
			)) as percentage
		FROM incident_reports ir
		WHERE ir.is_active = true
		AND ST_DistanceSphere(
			ST_MakePoint(ir.longitude, ir.latitude),
			ST_MakePoint($4, $5)
		) <= $6
		GROUP BY ir.category_code
		ORDER BY count DESC
		LIMIT 5
	`

	rows, err := ba.db.Query(query, lon, lat, radiusMeters, lon, lat, radiusMeters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryStats
	for rows.Next() {
		var cat CategoryStats
		err := rows.Scan(&cat.Category, &cat.Count, &cat.Percentage)
		if err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

// assessLocationDataQuality assesses data quality for a specific location
func (ba *BasicAnalytics) assessLocationDataQuality(incidents, activeUsers int) (LocationDataQuality, error) {
	quality := LocationDataQuality{}

	// Calculate data completeness for location
	completeness := 0.0
	if incidents >= 20 {
		completeness += 40
	}
	if incidents >= 50 {
		completeness += 30
	}
	if activeUsers >= 5 {
		completeness += 30
	}

	quality.DataCompleteness = completeness
	quality.HasMinimumData = completeness >= 50

	// Generate location-specific recommendations
	if incidents < 20 {
		quality.Recommendations = append(quality.Recommendations,
			"Need at least 20 incidents in this area for basic AI predictions")
	}
	if incidents < 50 {
		quality.Recommendations = append(quality.Recommendations,
			"Need at least 50 incidents in this area for advanced AI predictions")
	}
	if activeUsers < 5 {
		quality.Recommendations = append(quality.Recommendations,
			"Need at least 5 active users in this area for pattern recognition")
	}

	// Determine risk level based on incident density
	if incidents >= 50 {
		quality.RiskLevel = "high"
	} else if incidents >= 20 {
		quality.RiskLevel = "medium"
	} else {
		quality.RiskLevel = "low"
	}

	return quality, nil
}
