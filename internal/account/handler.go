package account

import (
	"alertly/internal/auth"
	"alertly/internal/database"
	"alertly/internal/response"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Original Functions ---

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

// --- New Apple IAP Validation Logic ---

// AppleReceiptRequest is the request body sent from the client
type AppleReceiptRequest struct {
	ReceiptData string `json:"receipt_data" binding:"required"`
	Platform    string `json:"platform"` // e.g., "ios"
}

// AppleVerifyRequest is the request sent to Apple's verifyReceipt endpoint
type AppleVerifyRequest struct {
	ReceiptData string `json:"receipt-data" binding:"required"`
	Password    string `json:"password,omitempty"` // Your app-specific shared secret
}

// AppleVerifyResponse is the top-level response from Apple
type AppleVerifyResponse struct {
	Status             int                      `json:"status"`
	Environment        string                   `json:"environment"`
	Receipt            AppleReceipt             `json:"receipt"`
	LatestReceiptInfo  []AppleLatestReceiptInfo `json:"latest_receipt_info"`
	LatestReceipt      string                   `json:"latest_receipt"`
	PendingRenewalInfo []PendingRenewalInfo     `json:"pending_renewal_info"`
}

// AppleReceipt contains general receipt information
type AppleReceipt struct {
	ReceiptType string `json:"receipt_type"`
	// ... other fields if needed
}

// AppleLatestReceiptInfo contains details about a specific transaction
type AppleLatestReceiptInfo struct {
	ProductID       string `json:"product_id"`
	TransactionID   string `json:"transaction_id"`
	ExpiresDateMS   string `json:"expires_date_ms"`
	IsTrialPeriod   string `json:"is_trial_period"`
	OriginalPurchaseDateMS string `json:"original_purchase_date_ms"`
}

// PendingRenewalInfo contains info about pending renewals
type PendingRenewalInfo struct {
	AutoRenewStatus string `json:"auto_renew_status"`
	ProductID       string `json:"product_id"`
}

const (
	appleProductionURL = "https://buy.itunes.apple.com/verifyReceipt"
	appleSandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
)

// ValidateAppleReceipt handles the validation of an Apple receipt
func ValidateAppleReceipt(c *gin.Context) {
	var accountID int64
	var err error
	var req AppleReceiptRequest

	// 1. Get user from JWT token
	accountID, err = auth.GetUserFromContext(c)
	if err != nil {
		log.Printf("Error getting user from context: %v", err)
		response.Send(c, http.StatusUnauthorized, true, "Unauthorized", nil)
		return
	}

	// 2. Bind JSON request from client
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid request body", nil)
		return
	}

	// 3. Validate the receipt with Apple
	appleResponse, err := verifyReceipt(req.ReceiptData)
	if err != nil {
		log.Printf("Error verifying receipt for account %d: %v", accountID, err)
		response.Send(c, http.StatusInternalServerError, true, "Error verifying receipt", nil)
		return
	}

	// 4. Check for a valid subscription
	if len(appleResponse.LatestReceiptInfo) == 0 {
		log.Printf("No active subscriptions found for account %d", accountID)
		response.Send(c, http.StatusBadRequest, true, "No active subscriptions found", nil)
		return
	}

	// Find the latest expiration date from all subscriptions
	var latestExpiration time.Time
	var productID string
	var hasValidSubscription bool

	for _, info := range appleResponse.LatestReceiptInfo {
		expiresMillis, err := strconv.ParseInt(info.ExpiresDateMS, 10, 64)
		if err != nil {
			continue // Skip if date is invalid
		}
		expirationDate := time.Unix(0, expiresMillis*int64(time.Millisecond))

		// Check if the subscription is currently active
		if expirationDate.After(time.Now()) {
			hasValidSubscription = true
			if expirationDate.After(latestExpiration) {
				latestExpiration = expirationDate
				productID = info.ProductID
			}
		}
	}

	if !hasValidSubscription {
		log.Printf("Subscription expired for account %d", accountID)
		response.Send(c, http.StatusBadRequest, true, "Subscription is expired", nil)
		return
	}

	// 5. Update premium status in the database
	repo := NewRepository(database.DB)
	service := NewService(repo)

	err = service.UpdatePremiumStatus(accountID, true, productID, &latestExpiration, req.Platform)
	if err != nil {
		log.Printf("Error updating premium status for account %d: %v", accountID, err)
		response.Send(c, http.StatusInternalServerError, true, "Error updating premium status", nil)
		return
	}

	log.Printf("✅ Premium status validated and updated for account %d. Product: %s, Expires: %s",
		accountID, productID, latestExpiration.String())

	response.Send(c, http.StatusOK, false, "Premium status updated successfully", gin.H{
		"product_id": productID,
		"expires_at": latestExpiration,
	})
}

// verifyReceipt sends the receipt data to Apple and handles sandbox retries
func verifyReceipt(receiptData string) (*AppleVerifyResponse, error) {
	// Start with production URL
	appleURL := appleProductionURL

	// Create the request to Apple
	appleReq := AppleVerifyRequest{
		ReceiptData: receiptData,
		// TODO: Add your App-Specific Shared Secret here if you have one
		// Password:    "your-shared-secret",
	}

	for i := 0; i < 2; i++ { // Allow one retry for sandbox
		reqBytes, err := json.Marshal(appleReq)
		if err != nil {
			return nil, err
		}

		// Make the POST request
		resp, err := http.Post(appleURL, "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var appleResp AppleVerifyResponse
		if err := json.Unmarshal(body, &appleResp); err != nil {
			return nil, err
		}

		// If status is 21007, it's a sandbox receipt. Retry with the sandbox URL.
		if appleResp.Status == 21007 && i == 0 {
			log.Println("Sandbox receipt detected, retrying with sandbox URL...")
			appleURL = appleSandboxURL
			continue
		}

		// If status is not 0, there's an error with the receipt
		if appleResp.Status != 0 {
			log.Printf("Invalid receipt. Apple status code: %d", appleResp.Status)
			return nil, &AppleValidationError{StatusCode: appleResp.Status}
		}

		return &appleResp, nil
	}

	return nil, &AppleValidationError{StatusCode: -1, Message: "Failed to validate receipt after retry."}
}

// AppleValidationError is a custom error for Apple validation failures
type AppleValidationError struct {
	StatusCode int
	Message    string
}

func (e *AppleValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "apple receipt validation failed with status " + strconv.Itoa(e.StatusCode)
}