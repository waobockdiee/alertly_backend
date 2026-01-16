package ai

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// EmergencyChatbot provides AI-powered emergency assistance
type EmergencyChatbot struct {
	db *sql.DB
}

// ChatMessage represents a conversation message
type ChatMessage struct {
	MessageID   int64     `json:"message_id"`
	UserID      int64     `json:"user_id"`
	Message     string    `json:"message"`
	IsUser      bool      `json:"is_user"`
	Timestamp   time.Time `json:"timestamp"`
	Category    string    `json:"category"`
	Priority    string    `json:"priority"`
}

// EmergencyResponse represents an AI-generated emergency response
type EmergencyResponse struct {
	Response    string   `json:"response"`
	Actions     []string `json:"actions"`
	Priority    string   `json:"priority"`
	Category    string   `json:"category"`
	Resources   []string `json:"resources"`
	NextSteps   []string `json:"next_steps"`
}

// NewEmergencyChatbot creates a new emergency chatbot
func NewEmergencyChatbot(db *sql.DB) *EmergencyChatbot {
	return &EmergencyChatbot{db: db}
}

// ProcessMessage processes a user message and returns an appropriate response
func (ec *EmergencyChatbot) ProcessMessage(userID int64, message string, context map[string]interface{}) (*EmergencyResponse, error) {
	// Analyze the message for emergency keywords and intent
	analysis := ec.analyzeMessage(message)
	
	// Generate appropriate response based on analysis
	response := ec.generateResponse(analysis, context)
	
	// Store the conversation
	err := ec.storeConversation(userID, message, response.Response, analysis.Category, analysis.Priority)
	if err != nil {
		return nil, fmt.Errorf("failed to store conversation: %w", err)
	}
	
	return response, nil
}

// analyzeMessage analyzes the user's message for emergency content
func (ec *EmergencyChatbot) analyzeMessage(message string) MessageAnalysis {
	message = strings.ToLower(message)
	
	analysis := MessageAnalysis{
		Category: "general",
		Priority: "low",
		Keywords: []string{},
		Intent:   "information",
	}
	
	// Emergency keywords detection
	emergencyKeywords := map[string]string{
		"emergency": "emergency",
		"help": "emergency",
		"urgent": "emergency",
		"danger": "emergency",
		"accident": "accident",
		"crash": "accident",
		"injury": "medical",
		"hurt": "medical",
		"bleeding": "medical",
		"fire": "fire",
		"smoke": "fire",
		"burning": "fire",
		"crime": "crime",
		"robbery": "crime",
		"theft": "crime",
		"assault": "crime",
		"attack": "crime",
		"suspicious": "suspicious",
		"strange": "suspicious",
		"weird": "suspicious",
	}
	
	// Check for emergency keywords
	for keyword, category := range emergencyKeywords {
		if strings.Contains(message, keyword) {
			analysis.Category = category
			analysis.Priority = "high"
			analysis.Keywords = append(analysis.Keywords, keyword)
		}
	}
	
	// Determine intent
	if strings.Contains(message, "what") || strings.Contains(message, "how") || strings.Contains(message, "where") {
		analysis.Intent = "question"
	} else if strings.Contains(message, "report") || strings.Contains(message, "submit") {
		analysis.Intent = "report"
	} else if analysis.Priority == "high" {
		analysis.Intent = "emergency"
	}
	
	return analysis
}

// generateResponse generates an appropriate response based on message analysis
func (ec *EmergencyChatbot) generateResponse(analysis MessageAnalysis, context map[string]interface{}) *EmergencyResponse {
	response := &EmergencyResponse{
		Priority: analysis.Priority,
		Category: analysis.Category,
	}
	
	switch analysis.Category {
	case "emergency":
		response.Response = "🚨 This sounds like an emergency. Please call 911 immediately if you're in immediate danger. I can help guide you through the situation while you wait for emergency services."
		response.Actions = []string{
			"Call 911 if immediate danger",
			"Move to a safe location",
			"Stay on the line with emergency services",
		}
		response.Resources = []string{
			"Emergency Services: 911",
			"Local Police: [Local Number]",
			"Fire Department: [Local Number]",
		}
		
	case "accident":
		response.Response = "🚗 I understand you're reporting an accident. Let me help you with the next steps."
		response.Actions = []string{
			"Ensure everyone is safe",
			"Call emergency services if needed",
			"Exchange information with other parties",
			"Take photos of the scene",
			"Report to insurance company",
		}
		response.NextSteps = []string{
			"Use Alertly to report the incident",
			"Contact your insurance provider",
			"Seek medical attention if injured",
		}
		
	case "medical":
		response.Response = "🏥 I see you're dealing with a medical situation. Here's what you should do:"
		response.Actions = []string{
			"Call 911 for serious injuries",
			"Apply first aid if trained",
			"Keep the person calm and still",
			"Monitor vital signs",
		}
		response.Resources = []string{
			"Emergency Medical Services: 911",
			"Poison Control: 1-800-222-1222",
			"Local Hospital: [Local Number]",
		}
		
	case "fire":
		response.Response = "🔥 Fire emergency detected. Please follow these safety protocols:"
		response.Actions = []string{
			"Get out immediately",
			"Call 911",
			"Don't use elevators",
			"Stay low to avoid smoke",
			"Meet at designated meeting point",
		}
		response.Resources = []string{
			"Fire Department: 911",
			"Emergency Services: 911",
		}
		
	case "crime":
		response.Response = "🚔 I understand you're reporting a crime. Here's what you should do:"
		response.Actions = []string{
			"Call 911 if crime is in progress",
			"Stay safe and don't confront",
			"Preserve evidence",
			"Get to a safe location",
		}
		response.NextSteps = []string{
			"Report to local police",
			"Use Alertly to document the incident",
			"Contact your insurance if property damage",
		}
		
	case "suspicious":
		response.Response = "👀 I see you've noticed something suspicious. Here's how to handle it:"
		response.Actions = []string{
			"Don't approach suspicious individuals",
			"Note descriptions and details",
			"Call non-emergency police line",
			"Stay in well-lit areas",
		}
		response.NextSteps = []string{
			"Report through Alertly",
			"Contact local police if needed",
		}
		
	default:
		response.Response = "I'm here to help! How can I assist you with your safety concerns today?"
		response.Actions = []string{
			"Report an incident through Alertly",
			"Get safety tips for your area",
			"Learn about emergency procedures",
		}
	}
	
	// Add context-specific information
	if location, exists := context["location"]; exists {
		response.Response += fmt.Sprintf("\n\n📍 I see you're in the %s area. I can provide location-specific guidance if needed.", location)
	}
	
	return response
}

// storeConversation stores the conversation in the database
func (ec *EmergencyChatbot) storeConversation(userID int64, userMessage, botResponse, category, priority string) error {
	query := `
		INSERT INTO emergency_chatbot_conversations
		(user_id, user_message, bot_response, category, priority, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := ec.db.Exec(query, userID, userMessage, botResponse, category, priority)
	return err
}

// GetConversationHistory retrieves conversation history for a user
func (ec *EmergencyChatbot) GetConversationHistory(userID int64, limit int) ([]ChatMessage, error) {
	query := `
		SELECT
			message_id,
			user_id,
			user_message as message,
			'user' as is_user,
			created_at as timestamp,
			category,
			priority
		FROM emergency_chatbot_conversations
		WHERE user_id = $1
		UNION ALL
		SELECT
			message_id,
			user_id,
			bot_response as message,
			'bot' as is_user,
			created_at as timestamp,
			category,
			priority
		FROM emergency_chatbot_conversations
		WHERE user_id = $2
		ORDER BY timestamp DESC
		LIMIT $3
	`

	rows, err := ec.db.Query(query, userID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		var isUserStr string
		err := rows.Scan(
			&msg.MessageID,
			&msg.UserID,
			&msg.Message,
			&isUserStr,
			&msg.Timestamp,
			&msg.Category,
			&msg.Priority,
		)
		if err != nil {
			continue
		}
		msg.IsUser = isUserStr == "user"
		messages = append(messages, msg)
	}
	
	return messages, nil
}

// MessageAnalysis represents the analysis of a user message
type MessageAnalysis struct {
	Category string   `json:"category"`
	Priority string   `json:"priority"`
	Keywords []string `json:"keywords"`
	Intent   string   `json:"intent"`
}
