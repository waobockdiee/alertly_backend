package ai

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

// PredictiveAnalytics provides AI-powered insights and predictions
type PredictiveAnalytics struct {
	db *sql.DB
}

// IncidentPrediction predicts likelihood of incidents in specific areas
type IncidentPrediction struct {
	AreaID             string              `json:"area_id"`
	Latitude           float64             `json:"latitude"`
	Longitude          float64             `json:"longitude"`
	RiskLevel          string              `json:"risk_level"` // low, medium, high, critical
	RiskScore          float64             `json:"risk_score"` // 0-100
	PredictedIncidents []PredictedIncident `json:"predicted_incidents"`
	Confidence         float64             `json:"confidence"`
	LastUpdated        time.Time           `json:"last_updated"`
}

type PredictedIncident struct {
	Category    string  `json:"category"`
	Probability float64 `json:"probability"` // 0-1
	TimeWindow  string  `json:"time_window"` // "next_hour", "next_day", "next_week"
	Severity    string  `json:"severity"`    // low, medium, high
}

// NewPredictiveAnalytics creates a new predictive analytics service
func NewPredictiveAnalytics(db *sql.DB) *PredictiveAnalytics {
	return &PredictiveAnalytics{db: db}
}

// AnalyzeAreaRisk analyzes risk patterns in a specific area
func (pa *PredictiveAnalytics) AnalyzeAreaRisk(latitude, longitude float64, radiusMeters int) (*IncidentPrediction, error) {
	// Get historical data for the area
	historicalData, err := pa.getHistoricalData(latitude, longitude, radiusMeters, 30) // 30 days
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data: %w", err)
	}

	// Calculate risk factors
	riskFactors := pa.calculateRiskFactors(historicalData)

	// Generate predictions
	predictions := pa.generatePredictions(historicalData, riskFactors)

	// Calculate overall risk score
	riskScore := pa.calculateOverallRiskScore(riskFactors)
	riskLevel := pa.determineRiskLevel(riskScore)

	return &IncidentPrediction{
		AreaID:             fmt.Sprintf("%.6f,%.6f", latitude, longitude),
		Latitude:           latitude,
		Longitude:          longitude,
		RiskLevel:          riskLevel,
		RiskScore:          riskScore,
		PredictedIncidents: predictions,
		Confidence:         pa.calculateConfidence(historicalData),
		LastUpdated:        time.Now(),
	}, nil
}

// getHistoricalData retrieves historical incident data for analysis
func (pa *PredictiveAnalytics) getHistoricalData(lat, lon float64, radiusMeters int, daysBack int) ([]HistoricalIncident, error) {
	query := `
		SELECT 
			ic.category_code,
			ic.subcategory_code,
			ic.created_at,
			ic.center_latitude,
			ic.center_longitude,
			ic.counter_total_votes_true,
			ic.counter_total_votes_false,
			COUNT(ir.inre_id) as incident_count
		FROM incident_clusters ic
		LEFT JOIN incident_reports ir ON ic.incl_id = ir.incl_id
		WHERE 
			ST_Distance_Sphere(
				POINT(ic.center_longitude, ic.center_latitude),
				POINT(?, ?)
			) <= ?
			AND ic.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY ic.incl_id, ic.category_code, ic.subcategory_code, ic.created_at, ic.center_latitude, ic.center_longitude, ic.counter_total_votes_true, ic.counter_total_votes_false
		ORDER BY ic.created_at DESC
	`

	rows, err := pa.db.Query(query, lon, lat, radiusMeters, daysBack)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []HistoricalIncident
	for rows.Next() {
		var incident HistoricalIncident
		err := rows.Scan(
			&incident.CategoryCode,
			&incident.SubcategoryCode,
			&incident.CreatedAt,
			&incident.Latitude,
			&incident.Longitude,
			&incident.VotesTrue,
			&incident.VotesFalse,
			&incident.Count,
		)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}

// calculateRiskFactors analyzes patterns to determine risk factors
func (pa *PredictiveAnalytics) calculateRiskFactors(incidents []HistoricalIncident) map[string]float64 {
	factors := make(map[string]float64)

	// Time-based patterns
	hourlyPatterns := make(map[int]int)
	weeklyPatterns := make(map[int]int)

	for _, incident := range incidents {
		hour := incident.CreatedAt.Hour()
		weekday := int(incident.CreatedAt.Weekday())

		hourlyPatterns[hour]++
		weeklyPatterns[weekday]++
	}

	// Calculate time-based risk
	factors["time_risk"] = pa.calculateTimeRisk(hourlyPatterns, weeklyPatterns)

	// Category-based risk
	categoryRisks := pa.calculateCategoryRisk(incidents)
	for category, risk := range categoryRisks {
		factors[category] = risk
	}

	// Density risk (incidents per area)
	factors["density_risk"] = pa.calculateDensityRisk(incidents)

	// Credibility risk
	factors["credibility_risk"] = pa.calculateCredibilityRisk(incidents)

	return factors
}

// generatePredictions creates AI-powered predictions
func (pa *PredictiveAnalytics) generatePredictions(incidents []HistoricalIncident, riskFactors map[string]float64) []PredictedIncident {
	var predictions []PredictedIncident

	// Analyze patterns for each category
	categories := pa.getUniqueCategories(incidents)

	for _, category := range categories {
		categoryIncidents := pa.filterByCategory(incidents, category)

		if len(categoryIncidents) > 0 {
			// Calculate probability based on historical frequency and current risk factors
			probability := pa.calculateCategoryProbability(categoryIncidents, riskFactors)

			if probability > 0.1 { // Only predict if probability > 10%
				prediction := PredictedIncident{
					Category:    category,
					Probability: probability,
					TimeWindow:  pa.determineTimeWindow(categoryIncidents),
					Severity:    pa.determineSeverity(categoryIncidents, riskFactors),
				}
				predictions = append(predictions, prediction)
			}
		}
	}

	return predictions
}

// calculateOverallRiskScore computes the overall risk score (0-100)
func (pa *PredictiveAnalytics) calculateOverallRiskScore(factors map[string]float64) float64 {
	// Weighted average of all risk factors
	weights := map[string]float64{
		"time_risk":         0.2,
		"density_risk":      0.3,
		"credibility_risk":  0.2,
		"crime":             0.15,
		"traffic_accident":  0.1,
		"medical_emergency": 0.05,
	}

	totalScore := 0.0
	totalWeight := 0.0

	for factor, weight := range weights {
		if score, exists := factors[factor]; exists {
			totalScore += score * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0
	}

	return math.Min(100, (totalScore/totalWeight)*100)
}

// determineRiskLevel converts risk score to human-readable level
func (pa *PredictiveAnalytics) determineRiskLevel(score float64) string {
	switch {
	case score < 25:
		return "low"
	case score < 50:
		return "medium"
	case score < 75:
		return "high"
	default:
		return "critical"
	}
}

// Helper methods
type HistoricalIncident struct {
	CategoryCode    string    `json:"category_code"`
	SubcategoryCode string    `json:"subcategory_code"`
	CreatedAt       time.Time `json:"created_at"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	VotesTrue       int       `json:"votes_true"`
	VotesFalse      int       `json:"votes_false"`
	Count           int       `json:"count"`
}

// calculateTimeRisk calculates real time-based risk patterns
func (pa *PredictiveAnalytics) calculateTimeRisk(hourly, weekly map[int]int) float64 {
	if len(hourly) == 0 || len(weekly) == 0 {
		return 0.3 // Default risk if no data
	}

	// Analyze peak hours (18:00-22:00 typically have more incidents)
	peakHours := []int{18, 19, 20, 21, 22}
	peakHourIncidents := 0
	totalIncidents := 0

	for hour, count := range hourly {
		totalIncidents += count
		for _, peakHour := range peakHours {
			if hour == peakHour {
				peakHourIncidents += count
			}
		}
	}

	// Calculate peak hour risk (0-1)
	peakHourRisk := 0.0
	if totalIncidents > 0 {
		peakHourRisk = float64(peakHourIncidents) / float64(totalIncidents)
	}

	// Analyze weekend patterns (Friday-Sunday typically have more incidents)
	weekendDays := []int{5, 6, 0} // Friday, Saturday, Sunday
	weekendIncidents := 0
	totalWeeklyIncidents := 0

	for day, count := range weekly {
		totalWeeklyIncidents += count
		for _, weekendDay := range weekendDays {
			if day == weekendDay {
				weekendIncidents += count
			}
		}
	}

	// Calculate weekend risk (0-1)
	weekendRisk := 0.0
	if totalWeeklyIncidents > 0 {
		weekendRisk = float64(weekendIncidents) / float64(totalWeeklyIncidents)
	}

	// Combine risks (weighted average)
	timeRisk := (peakHourRisk * 0.6) + (weekendRisk * 0.4)
	return math.Min(1.0, timeRisk)
}

func (pa *PredictiveAnalytics) calculateCategoryRisk(incidents []HistoricalIncident) map[string]float64 {
	risks := make(map[string]float64)
	categoryCounts := make(map[string]int)

	for _, incident := range incidents {
		categoryCounts[incident.CategoryCode] += incident.Count
	}

	total := 0
	for _, count := range categoryCounts {
		total += count
	}

	for category, count := range categoryCounts {
		risks[category] = float64(count) / float64(total)
	}

	return risks
}

func (pa *PredictiveAnalytics) calculateDensityRisk(incidents []HistoricalIncident) float64 {
	// Calculate incidents per square kilometer
	return float64(len(incidents)) / 100.0 // Placeholder
}

func (pa *PredictiveAnalytics) calculateCredibilityRisk(incidents []HistoricalIncident) float64 {
	totalVotes := 0
	trueVotes := 0

	for _, incident := range incidents {
		totalVotes += incident.VotesTrue + incident.VotesFalse
		trueVotes += incident.VotesTrue
	}

	if totalVotes == 0 {
		return 0.5
	}

	// Lower credibility = higher risk
	return 1.0 - (float64(trueVotes) / float64(totalVotes))
}

func (pa *PredictiveAnalytics) getUniqueCategories(incidents []HistoricalIncident) []string {
	categories := make(map[string]bool)
	for _, incident := range incidents {
		categories[incident.CategoryCode] = true
	}

	var result []string
	for category := range categories {
		result = append(result, category)
	}
	return result
}

func (pa *PredictiveAnalytics) filterByCategory(incidents []HistoricalIncident, category string) []HistoricalIncident {
	var filtered []HistoricalIncident
	for _, incident := range incidents {
		if incident.CategoryCode == category {
			filtered = append(filtered, incident)
		}
	}
	return filtered
}

// calculateCategoryProbability calculates real probability based on historical patterns
func (pa *PredictiveAnalytics) calculateCategoryProbability(incidents []HistoricalIncident, riskFactors map[string]float64) float64 {
	if len(incidents) == 0 {
		return 0.1 // Default low probability
	}

	// Calculate base probability from historical frequency
	totalIncidents := len(incidents)
	categoryIncidents := len(pa.filterByCategory(incidents, incidents[0].CategoryCode))
	baseProbability := float64(categoryIncidents) / float64(totalIncidents)

	// Apply time-based risk factor
	timeRisk := riskFactors["time_risk"]
	if timeRisk > 0 {
		baseProbability *= (1 + timeRisk)
	}

	// Apply density risk factor
	densityRisk := riskFactors["density_risk"]
	if densityRisk > 0 {
		baseProbability *= (1 + densityRisk*0.5)
	}

	// Apply credibility risk factor (lower credibility = higher risk)
	credibilityRisk := riskFactors["credibility_risk"]
	if credibilityRisk > 0 {
		baseProbability *= (1 + credibilityRisk*0.3)
	}

	// Normalize probability (0-1)
	probability := math.Min(1.0, baseProbability)

	// Apply confidence based on data quality
	confidence := pa.calculateConfidence(incidents)
	probability *= confidence

	return probability
}

// determineTimeWindow analyzes patterns to determine most likely time window
func (pa *PredictiveAnalytics) determineTimeWindow(incidents []HistoricalIncident) string {
	if len(incidents) < 3 {
		return "next_day" // Default if insufficient data
	}

	// Analyze recent activity patterns
	recentIncidents := 0
	totalIncidents := len(incidents)

	// Count incidents in last 7 days
	weekAgo := time.Now().AddDate(0, 0, -7)
	for _, incident := range incidents {
		if incident.CreatedAt.After(weekAgo) {
			recentIncidents++
		}
	}

	// Calculate activity rate
	activityRate := float64(recentIncidents) / float64(totalIncidents)

	// Determine time window based on activity rate
	if activityRate > 0.5 {
		return "next_hour" // High activity = immediate risk
	} else if activityRate > 0.2 {
		return "next_day" // Moderate activity = daily risk
	} else {
		return "next_week" // Low activity = weekly risk
	}
}

// determineSeverity analyzes incident severity patterns
func (pa *PredictiveAnalytics) determineSeverity(incidents []HistoricalIncident, riskFactors map[string]float64) string {
	if len(incidents) == 0 {
		return "low"
	}

	// Calculate severity based on multiple factors
	severityScore := 0.0

	// Factor 1: Incident density
	densityRisk := riskFactors["density_risk"]
	severityScore += densityRisk * 0.4

	// Factor 2: Credibility risk (lower credibility = higher severity)
	credibilityRisk := riskFactors["credibility_risk"]
	severityScore += credibilityRisk * 0.3

	// Factor 3: Time risk
	timeRisk := riskFactors["time_risk"]
	severityScore += timeRisk * 0.3

	// Determine severity level
	if severityScore < 0.3 {
		return "low"
	} else if severityScore < 0.6 {
		return "medium"
	} else {
		return "high"
	}
}

func (pa *PredictiveAnalytics) calculateConfidence(incidents []HistoricalIncident) float64 {
	// Calculate confidence based on data quality and quantity
	if len(incidents) < 5 {
		return 0.3
	} else if len(incidents) < 20 {
		return 0.6
	} else {
		return 0.9
	}
}
