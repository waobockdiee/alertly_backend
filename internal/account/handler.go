package account

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMyInfo(c *gin.Context) {
	var accountID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	ip := c.ClientIP()

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetMyInfo(accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting history", nil)
		return
	}

	// Save the last request for the account for cronjob method(send notification push to user)
	err = service.SaveLastRequest(accountID, ip)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	response.Send(c, http.StatusOK, false, "success", data)

}

func GetHistory(c *gin.Context) {
	var accountID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetHistory(accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting history", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)

}

func ClearHistory(c *gin.Context) {
	var accountID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.ClearHistory(accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func DeleteAccount(c *gin.Context) {
	response.Send(c, http.StatusOK, false, "success", nil)
}

func GetCounterHistories(c *gin.Context) {
	var accountID int64
	var err error
	var counter Counter

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	counter, err = service.GetCounterHistories(accountID)

	log.Printf("COUNTER: %v", counter)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", counter)
}

func SetHasFinishedTutorial(c *gin.Context) {
	var accountID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.SetHasFinishedTutorial(accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", nil)
}

func GetViewedIncidentIds(c *gin.Context) {
	var accountID int64
	var err error

	accountID, err = auth.GetUserFromContext(c)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	data, err := service.GetViewedIncidentIds(accountID)

	if err != nil {
		log.Printf("Error: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error getting viewed incident IDs", nil)
		return
	}

	response.Send(c, http.StatusOK, false, "success", data)
}

// UpdatePremiumStatusRequest represents the request body for updating premium status
type UpdatePremiumStatusRequest struct {
	IsPremium        bool   `json:"is_premium" binding:"required"`
	SubscriptionType string `json:"subscription_type"`
	PurchaseDate     string `json:"purchase_date"`
	Platform         string `json:"platform"`
}

// UpdatePremiumStatus handles updating the user's premium subscription status
func UpdatePremiumStatus(c *gin.Context) {
	var accountID int64
	var err error
	var req UpdatePremiumStatusRequest

	// Get user from JWT token
	accountID, err = auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid request body", nil)
		return
	}

	// Validate purchase date format if provided
	var purchaseDate *time.Time
	if req.PurchaseDate != "" {
		parsedDate, err := time.Parse(time.RFC3339, req.PurchaseDate)
		if err != nil {
			log.Printf("Error parsing purchase date: %v", err)
			response.Send(c, http.StatusBadRequest, true, "Invalid purchase date format", nil)
			return
		}
		purchaseDate = &parsedDate
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	// Update premium status in database
	err = service.UpdatePremiumStatus(accountID, req.IsPremium, req.SubscriptionType, purchaseDate, req.Platform)
	if err != nil {
		log.Printf("Error updating premium status: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "Error updating premium status", nil)
		return
	}

	log.Printf("✅ Premium status updated for account %d: isPremium=%v, type=%s, platform=%s", 
		accountID, req.IsPremium, req.SubscriptionType, req.Platform)

	response.Send(c, http.StatusOK, false, "Premium status updated successfully", nil)
}
