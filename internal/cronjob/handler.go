package cronjob

import (
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RunPremiumExpirationCheck manually triggers the premium expiration check
func RunPremiumExpirationCheck(c *gin.Context) {
	log.Println("🔄 Manual premium expiration check triggered")

	service := NewPremiumExpirationService(database.DB)

	err := service.CheckAndExpirePremiumAccounts()
	if err != nil {
		log.Printf("❌ Premium expiration check failed: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Premium expiration check failed", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Premium expiration check completed successfully", nil)
}

// GetPremiumStats returns statistics about premium subscriptions
func GetPremiumStats(c *gin.Context) {
	service := NewPremiumExpirationService(database.DB)

	stats, err := service.GetPremiumExpirationStats()
	if err != nil {
		log.Printf("❌ Error getting premium stats: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting premium statistics", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Premium statistics retrieved successfully", stats)
}

// SendExpirationWarnings manually triggers expiration warnings
func SendExpirationWarnings(c *gin.Context) {
	log.Println("📧 Manual expiration warnings triggered")

	service := NewPremiumExpirationService(database.DB)

	err := service.SendExpirationWarnings()
	if err != nil {
		log.Printf("❌ Expiration warnings failed: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Expiration warnings failed", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "Expiration warnings sent successfully", nil)
}
