package signup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"alertly/internal/database"
	"alertly/internal/emails"
	"alertly/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func RegisterUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Invalid input. Please check the data and try again.", nil)
		return
	}

	// Log de debug para ver qué datos se están recibiendo
	log.Printf("DEBUG: Datos recibidos en signup - Email: %s, FirstName: %s, LastName: %s, Password: %s (length: %d)",
		user.Email, user.FirstName, user.LastName,
		strings.Repeat("*", len(user.Password)), len(user.Password))

	if err := validate.Struct(user); err != nil {
		log.Printf("Error de validación: %v", err)
		response.Send(c, http.StatusBadRequest, true, "Some form fields are invalid. Please review and try again.", nil)
		return
	}

	repo := NewRepository(database.DB)
	service := NewService(repo)

	registeredUser, code, err := service.RegisterUser(user)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		response.Send(c, http.StatusInternalServerError, true, "We couldn't register your account. Please try again later.", nil)
		return
	}

	// Si hay código de referral, registrar la conversión
	if user.ReferralCode != "" {
		go registerReferralConversion(user.ReferralCode, registeredUser.ID)
	}

	emails.SendTemplate(user.Email, "Activate your Alertly account", "new_account_activation_code", map[string]string{
		"Name": user.FirstName,
		"Code": code,
	})
	response.Send(c, http.StatusOK, false, "Your account has been created successfully. Please check your email to activate it.", registeredUser)
}

// registerReferralConversion registra la conversión de referral de forma asíncrona
func registerReferralConversion(referralCode string, userID int64) {
	apiKey := os.Getenv("REFERRAL_API_KEY")
	if apiKey == "" {
		log.Println("⚠️ REFERRAL_API_KEY not set, skipping referral conversion")
		return
	}

	// Preparar request body
	conversionData := map[string]interface{}{
		"referral_code": referralCode,
		"user_id":       userID,
		"registered_at": time.Now().Format(time.RFC3339),
		"platform":      "iOS", // Default, podría venir del frontend
	}

	jsonData, err := json.Marshal(conversionData)
	if err != nil {
		log.Printf("❌ Error marshaling referral conversion: %v", err)
		return
	}

	// Hacer llamada HTTP interna
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	url := fmt.Sprintf("http://localhost:%s/api/v1/referral/conversion", serverPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Error creating referral conversion request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error registering referral conversion: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		log.Printf("✅ Referral conversion registered successfully for user %d with code %s", userID, referralCode)
	} else {
		log.Printf("⚠️ Failed to register referral conversion. Status: %d, User: %d, Code: %s", resp.StatusCode, userID, referralCode)
	}
}
