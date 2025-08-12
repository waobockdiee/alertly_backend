package analytics

import (
	"alertly/internal/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	analytics *BasicAnalytics
}

func NewHandler(analytics *BasicAnalytics) *Handler {
	return &Handler{
		analytics: analytics,
	}
}

// GetAnalyticsSummary returns basic analytics summary
func (h *Handler) GetAnalyticsSummary(c *gin.Context) {
	summary, err := h.analytics.GetAnalyticsSummary()
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Failed to fetch analytics data", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "Analytics data retrieved successfully", summary)
}

// GetSimplePredictions returns basic predictions based on current data
func (h *Handler) GetSimplePredictions(c *gin.Context) {
	predictions, err := h.analytics.GetSimplePredictions()
	if err != nil {
		response.Send(c, http.StatusInternalServerError, true, "Failed to generate predictions", err.Error())
		return
	}

	response.Send(c, http.StatusOK, false, "Predictions generated successfully", predictions)
}

// TestAnalytics is a test endpoint to verify the analytics system is working
func (h *Handler) TestAnalytics(c *gin.Context) {
	response.Send(c, http.StatusOK, false, "Analytics system is working", map[string]interface{}{
		"status":  "ok",
		"message": "Analytics system is operational",
	})
}

// GetLocationAnalytics returns analytics for a specific location
func (h *Handler) GetLocationAnalytics(c *gin.Context) {
	// Get location parameters from query string
	latStr := c.Query("latitude")
	lonStr := c.Query("longitude")
	radiusStr := c.Query("radius")

	if latStr == "" || lonStr == "" || radiusStr == "" {
		response.Send(c, http.StatusBadRequest, true, "Missing required parameters: latitude, longitude, radius", nil)
		return
	}

	// Parse parameters
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Invalid latitude parameter", nil)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Invalid longitude parameter", nil)
		return
	}

	radius, err := strconv.Atoi(radiusStr)
	if err != nil {
		response.Send(c, http.StatusBadRequest, true, "Invalid radius parameter", nil)
		return
	}

	// Get analytics for this location
	analytics, err := h.analytics.GetLocationAnalytics(lat, lon, radius)
	if err != nil {
		log.Printf("Error getting location analytics: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting location analytics", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Location analytics retrieved successfully", analytics)
}
